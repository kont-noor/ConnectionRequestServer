require 'sinatra'
require 'yaml'
require 'redis'

class ConnectionRequestServer < Sinatra::Base
  configure do
    yaml = File.read(File.dirname(__FILE__) + '/config/redis.yml')
    redis_settings = YAML.load(yaml)
    REDIS = Redis.new(:host => redis_settings['host'], :port => redis_settings['port'], :password => redis_settings['password'])
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
      return 400 if self.user_connected?(param[:activation_code], param[:device_id])
      return self.connect_user(param[:activation_code], param[:device_id]) ? 1 : 500
    end
  end

  def self.user_connected? activation_code, device_id
    #REDIS.hset(1, 1000, 100)
    account_info = REDIS.hget(activation_code.to_i/1000, activation_code)
    return false if account_info.nil?
    return account_info != device_id 
  end

  def self.connect_user activation_code, device_id
    begin
      REDIS.hset(activation_code.to_i/1000, activation_code, device_id)
      return true
    resque Exception => e
      #TODO: add exception to log
      return false
    end
  end

  def disconnect param
    REDIS.get('s')
  end

  def heartbeat param
  end
end
