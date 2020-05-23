# reseed-monitoring
A service for monitoring I2P reseed servers remotely. It asks for a reseed bundle
from all available sources(defined in config.json) and outputs information about
them.

The easiest way to run it is with docker installed. I just run this script on a daily basis using
cron:


		#! /usr/bin/env sh
		docker pull  eyedeekay/reseed-monitoring
		docker rm -f reseed-monitoring
		docker run -itd --name reseed-monitoring --restart=always -p 127.0.0.1:7672:7672 eyedeekay/reseed-monitoring

Then, set up a an I2PTunnel to it in the Hidden Services Manager.
