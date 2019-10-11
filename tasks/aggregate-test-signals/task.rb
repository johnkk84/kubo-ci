#!/usr/bin/env ruby
# frozen_string_literal: true

require 'terminal-table'
require 'rainbow'

linux_files = Dir.glob('gcs-*-shipables').flat_map { |d| Dir.glob(d + '/*shipable') }
linux_files.delete_if { |f| /windows/.match(f) }
windows_files = Dir.glob('gcs-*windows*-shipables').flat_map { |d| Dir.glob(d + '/*shipable') }

def find_overlap(files)
  overlap = File.read(files.first).split("\n")
  files.each do |f|
    overlap = File.read(f).split("\n") & overlap
    puts "After checking #{f} good versions are: #{overlap.last}"
  end
  overlap
end

puts
puts "Looking for highest common green build..."
puts "....linux"
linux_overlap = find_overlap(linux_files)

if linux_overlap.any?
  linux_release_sha, deployment_sha, linux_build_number = linux_overlap.last.split
  File.write(ENV['SLACK_MESSAGE_FILE'],
             "Ready to :ship: <https://github.com/cloudfoundry-incubator/kubo-release/tree/#{linux_release_sha}/|#{linux_release_sha}> <https://github.com/cloudfoundry-incubator/kubo-deployment/tree/#{deployment_sha}/|#{deployment_sha}> Build number is #{linux_build_number}")
else
  puts 'No good versions yet'
  File.write(ENV['SLACK_MESSAGE_FILE'], 'No shippable version found')
  exit 1
end

puts Rainbow("Good versions are: #{linux_release_sha}, #{deployment_sha}. Build number is #{linux_build_number}").green

puts
puts "....windows"
windows_overlap = find_overlap(windows_files)

if windows_overlap.any?
  windows_release_sha, windows_deployment_sha, windows_build_number = windows_overlap.last.split
end

puts Rainbow("Good versions are: #{windows_release_sha}, #{windows_deployment_sha}. Build number is #{windows_build_number}").green

puts
puts "Highest green build for each pipeline..."
rows = []
files = linux_files + windows_files
files.each do |f|
  builds = File.read(f).split("\n")
  pipeline = f.split("/")[1].split("-shipable")[0]
  linux_release_sha, deployment_sha, linux_build_number = builds.last.split
  rows << [pipeline, linux_build_number]
end
table = Terminal::Table.new :headings => ['Pipeline', 'Build Number'], :rows => rows
puts table
puts

if linux_build_number == windows_build_number
  File.write(ENV['SHIPABLE_VERSION_FILE'], [linux_release_sha, windows_release_sha, deployment_sha, linux_build_number])
  puts "#{ENV['SHIPABLE_VERSION_FILE']} contains..."
  `cat #{ENV['SHIPABLE_VERSION_FILE']}`
else
  puts Rainbow("linux_build_number #{linux_build_number} does not match the windows_build_number #{windows_build_number} ").red
end
