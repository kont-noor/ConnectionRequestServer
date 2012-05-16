require_relative '../main'
require 'rack/test'

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

    it "should response on connect request" do
      last_response.should be_ok
    end

    it "should be an xml" do
      last_response.header['Content-Type'].should =~ /text\/xml/
    end
  
    it "should be a standard form" do
      last_response.body.should =~ /\<connection_request_response\>\s*\<code\>.*?\<\/code\>\s*\<message\>.*?\<\/message\>\s*\<\/connection_request_response\>/ims
    end
  end
end

describe "permission to connect" do
  describe "return 1 (approved)" do
    pending "when approved the connection"
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
    pending "on unknown error"
  end
end
