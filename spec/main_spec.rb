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
  describe "return 1 (approved)" do
    it "when approved the connection" do
      pending "implement later"
    end
  end

  describe "return 401 (missing parameters)" do
    it "on parameters missing" do
      app.request_permission_to_connect.should == 401
    end

    it "when parameters hash is empty" do
      app.request_permission_to_connect({}).should == 401
    end

    it "when parameter is not a hash" do
      app.request_permission_to_connect(1).should == 401
    end

    it "when activation_code is missing" do
      app.request_permission_to_connect({:device_id => '300'}).should == 401
    end

    it "when activation_code is empty" do
      app.request_permission_to_connect({:device_id => '300', :activation_code => ''}).should == 401
    end

    it "when device_id is empty" do
      app.request_permission_to_connect({:device_id => '', :activation_code => '300'}).should == 401
    end

    it "when device_id is missing" do
      app.request_permission_to_connect({:activation_code => '300'}).should == 401
    end
  end

  describe "return 500" do
    it "on unknown error" do
      pending "what unknown error?"
    end
  end
end
