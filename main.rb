require 'sinatra'
require 'yaml'
require 'redis'

class ConnectionRequestServer < Sinatra::Base
  configure do
    REDIS = Redis.new(:host => 'localhost', :port => 6379, :password => nil)
  end

  post '/request_permission_to_connect' do
    response_code = ConnectionRequestServer.request_permission_to_connect params

    message =  get_messages
    erb :connection_response, :locals => {:code => response_code, :message => message[response_code]}, :content_type => 'text/xml'
  end

  post '/disconnect' do
    disconnect params
    'ok'
  end

  post '/heartbeat' do
    heartbeat params
    'ok'
  end

  error 400..510 do
    '404'
  end

  def get_messages
    messages_yaml = File.read(File.dirname(__FILE__) + '/config/messages.yml')
    YAML.load(messages_yaml)
  end

  def self.request_permission_to_connect(param={})
    param = {} unless param.is_a?(Hash)
    param = {:activation_code => nil, :device_id => nil}.merge(param)

    if (param[:activation_code].nil? or param[:device_id].nil? or param[:activation_code].empty? or param[:device_id].empty?)
      return 401
    else
      return 500
    end

  end

  def disconnect param
    REDIS.get('s')
  end

  def heartbeat param
  end
end
