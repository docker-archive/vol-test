#!/bin/ruby

require 'sinatra'

set :bind, '0.0.0.0'
set :port, 80

get '/' do
  return "Docker Volume Certification Container"
end

get '/resetfilecheck' do
  `rm /data/*`
  return "1"
end

get '/status' do
  return "OK"
end

get '/shutdown' do
  Process.kill('TERM', Process.pid)
end

get '/runfilecheck' do
  `/usr/bin/idempotent_filecheck.sh`
end

get '/textcheck' do
  if File.exist?("/data/textfile")
    line = File.open("/data/textfile").first.strip
    if line == "dockertext"
      return "1"
    else
      return "0"
    end
  else
    return "0"
  end
end

get '/bincheck' do
  out = `md5sum -c /data/binchecksum`.strip!
  if out == "/data/binaryfile: OK"
    return "1"
  else
    return "0"
  end

end
