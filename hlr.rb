# HLR lib
# generated from sys/securite/iptables

require 'benchmark'

def stripnb(host)
	return host.sub(/\d+$/,'')
end

class ReadFile
	def initialize
		@done_files = {}
	end
	def read_one_file(file_name, lines)
		@done_files[file_name] = true
		File.foreach(file_name) do |line|
			if line =~ /^ \s* include \s+ "(.*?)"/x
				lines += read_file($1)
			else
				lines << line.chomp
			end
		end
		lines
	end
	def read_file(fname)
		lines = []
		return read_one_file(fname, lines) if fname[0] == '/'
		["hosts", "hosts/media", "includes"].each do |p|
			name = "#{p}/#{fname}"
			[name, stripnb(name)].uniq.each do |file_name|
				next unless File.exist?(file_name)
				next if @done_files.key?(file_name)
				lines = read_one_file(file_name, lines)
			end
		end
		return lines
	end
end

tries = 10000
Benchmark.bm(20) do |x|
	x.report("ReadFile #{tries}") {
		tries.times {
			read_file = ReadFile.new
			lines = ReadFile.new.read_file("test")
		}
	}
end

# vim: set ts=2 sw=2 list:
