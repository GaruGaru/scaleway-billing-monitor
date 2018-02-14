

## Billing monitor for [Scaleway](https://cloud.scaleway.com/#/) written in Go


### swarm deploy example 
	
	  scaleway-billing-monitor:
		  image: garugaru/scaleway-billing-monitor
		  environment:
			- STATSD_HOST=127.0.0.1:8125
			- SCALEWAY_AUTH_TOKEN=xxxxxx-xxxx-xx-xxxxx
		  deploy:
			mode: replicated
			replicas: 1
			restart_policy:
			  condition: on-failure
			  delay: 5s
			  max_attempts: 3
			  
