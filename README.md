# User Segmentation System
## Introduction
This system is designed to efficiently manage and retrieve the number of users in different segments.

## Architecture
The architecture of this system is composed of several components:

- Redis: used as a cache to store the number of users in each segment.
- Endpoints:
  - Endpoint for archiving expired data: this endpoint is used to remove the expired data from the cache. A job in Kubernetes is defined to call this endpoint periodically.
  - Endpoint for updating cache: this endpoint is used to update the number of users in each segment in the cache. A job in Kubernetes is defined to call this endpoint periodically.
  - Endpoint for retrieving user count in a specific segment: this endpoint is used to retrieve the number of users in a specific segment from the cache.
  - Endpoint for storing new pairs of user and segment: this endpoint is used to store the new pairs of user and segment in the cache.
- Load Balancer
- Postgres

![](media/arch.jpg)



## Scalability
This system is suitable for scaling as it separates the storage and retrieval of data, allowing for horizontal scaling of the cache. Additionally, the use of a distributed task queue such as Kubernetes allows for parallel processing and improved performance.

By utilizing a cache, the system also reduces the load on the database and improves the response time for retrieving user count in a specific segment.

In case of increased load, the cache can be easily scaled horizontally by adding more nodes to the cache cluster. The use of Kubernetes also allows for seamless scaling of the system components.