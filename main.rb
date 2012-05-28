require 'sinatra'
require 'yaml'
require 'redis'
require 'json'

class ConnectionRequestServer < Sinatra::Base

  @error = nil
  CONNECTION_PERIOD = 10000000

  configure do
    yaml = File.read(File.dirname(__FILE__) + '/config/redis.yml')
    redis_settings = YAML.load(yaml)
    begin
      REDIS = Redis.new(:host => redis_settings['host'], :port => redis_settings['port'], :password => redis_settings['password'])
    rescue Exception => e
      #TODO: add exception to log
      @error = e
    end
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
    return 500 unless @error.nil?
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
    account_info = REDIS.hget("connections:#{self.get_hash_name(activation_code)}", activation_code)
    return false if account_info.nil?
    account_info = JSON.parse(account_info)
    return account_info['device_id'] != device_id || Time.now.to_i - account_info['connection_time'].to_i > CONNECTION_PERIOD
  end

  def self.connect_user activation_code, device_id
    begin
      REDIS.hset("connections:#{self.get_hash_name(activation_code)}", activation_code, {:device_id => device_id, :connection_time => Time.now.to_i}.to_json)
      return true
    rescue Exception => e
      #TODO: add exception to log
      return false
    end
  end

  def disconnect param
    unless param[:activation_code].nil?
      REDIS.hdel("connections:#{self.get_hash_name(param[:activation_code])}", param[:activation_code])
    end
  end

  def self.get_hash_name code
    code.to_i/1000
  end

  def heartbeat param
    #TODO: add connection time for heartbeat
  end
end
