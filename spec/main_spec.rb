require_relative '../main'
require 'rack/test'

set :environment, :test

def app
  ConnectionRequestServer
end

describe 'base response' do
  include Rack::Test::Methods

  it "should show 404 message" do
    get '/'
    last_response.body.should == '404'
  end

  it "should show 404 message" do
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

  it "should response on connect request" do
    post '/request_permission_to_connect'
    last_response.should be_ok
  end

  it "should be an xml" do
    post '/request_permission_to_connect'
    last_response.header['Content-Type'].should =~ /text\/xml/
  end
  
  it "should be a standard connection response" do
    post '/request_permission_to_connect'
    last_response.body.should =~ /\<connection_request_response\>\s*\<code\>.*?\<\/code\>\s*\<message\>.*?\<\/message\>\s*\<\/connection_request_response\>/ims
  end
end
