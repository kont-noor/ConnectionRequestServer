require_relative '../main'
require 'rack/test'
require 'nokogiri'

set :environment, :test

def app
  ConnectionRequestServer
end

describe 'base response' do
  include Rack::Test::Methods

  it "get root should return 404 message" do
    get '/'
    last_response.body.should == '404'
  end

  it "post to root should return 404 message" do
    post '/', :xml => 'request'
    last_response.body.should == '404'
  end

  it "should response on disconnect" do
    post '/disconnect'
    last_response.body.should == 'ok'
  end

  it "should response on heartbeat" do
    post '/heartbeat'
    last_response.body.should == 'ok'
  end

  describe 'permission to connect response' do
    before(:all) do
      post '/request_permission_to_connect'
    end

    let(:xml){ Nokogiri::XML(last_response.body)}

    it "should response on connect request" do
      last_response.should be_ok
    end

    it "content-type should be an xml" do
      last_response.header['Content-Type'].should =~ /text\/xml/
    end

    it "should be a valid xml" do
      pending "implement later"
    end

    it "should have message code" do
      xml.css('connection_request_response code').should_not be_empty
    end

    it "should have message text" do
      xml.css('connection_request_response message').should_not be_empty
    end
  end
end

describe "permission to connect" do
  include Rack::Test::Methods

  let(:redis) do
    yaml = File.read(File.dirname(__FILE__) + '/../config/redis.yml')
    redis_settings = YAML.load(yaml)
    Redis.new(:host => redis_settings['host'], :port => redis_settings['port'], :password => redis_settings['password'])
  end

  let(:device1) {{'device_id' => '1', 'activation_code' => '1000'}}
  let(:device2) {{'device_id' => '3', 'activation_code' => '1000'}}
  let(:xml){ Nokogiri::XML(last_response.body)}
  let(:message){ xml.at_css('connection_request_response code').content}

  def permission dev
    app.api_request_permission_to_connect(dev)
  end

  describe "return 1 (approved)" do
    before(:each) do
      redis.del('connections:1')
    end

    before(:each, :turn => :request) do
      post '/request_permission_to_connect', device1
    end

    describe "when approved the connection" do
      it "for method" do
        permission(device1).should == 1
      end

      it "for request", :turn => :request do
        message.should == '1'
      end
    end

    describe "when user connected from the same device" do
      it "for method" do
        permission(device1)
        permission(device1).should == 1
      end

      it "for request", :turn => :request do
        post '/request_permission_to_connect', device1
        message.should == '1'
      end
    end

    describe "when disconnected by uptime" do
      it "for method" do
        permission(device1)
        sleep app::CONNECTION_PERIOD * 2
        permission(device2).should == 1
      end

      it "for request", :turn => :request do
        sleep app::CONNECTION_PERIOD * 2
        post '/request_permission_to_connect', device2
        message.should == '1'
      end
    end

    describe "when disconnected by device" do
      it "for method" do
        permission(device1)
        app.api_disconnect(device1)
        permission(device2).should == 1
      end

      it "for request", :turn => :request do
        post '/disconnect', device1
        post '/request_permission_to_connect', device2
        message.should == '1'
      end
    end
  end

  describe "return 400 (user already connected)" do
    before(:each) do
      redis.del('connections:1')
    end

    before(:each, :turn => :method) do
      permission(device1)
    end

    before(:each, :turn => :request) do
      post '/request_permission_to_connect', device1
    end

    describe "when user already connected from another device" do
      it "for method", :turn => :method do
        permission(device2).should == 400
      end

      it "for request", :turn => :request do
        post '/request_permission_to_connect', device2
        message.should == '400'
      end
    end

    describe "when not disconnected by uptime" do
      it "for method", :turn => :method do
        sleep app::CONNECTION_PERIOD / 2
        permission(device2).should == 400
      end

      it "for request", :turn => :request do
        sleep app::CONNECTION_PERIOD / 2
        post '/request_permission_to_connect', device2
        message.should == '400'
      end
    end
  end

  describe "return 401 (missing parameters)" do
    invalid_requests = {
      "parameters missing" => nil,
      "parameters hash is empty" => {},
      "parameter is not a hash" => '1',
      "activation_code is missing" => {'device_id' => '300'},
      "activation_code is empty" => {'device_id' => '300', 'activation_code' => ''},
      "device_id is empty" => {'device_id' => '', 'activation_code' => '300'},
      "device_id is missing" => {'activation_code' => '300'}
    }

    invalid_requests.each_pair do |key, value|
      describe "when #{key}" do
        it "for method" do
          permission(value).should == 401
        end

        it "for request" do
          post '/request_permission_to_connect', value
          message.should == '401'
        end
      end
    end
  end

  describe "return 500" do
    it "on unknown error" do
      pending "what unknown error?"
    end
  end
end
