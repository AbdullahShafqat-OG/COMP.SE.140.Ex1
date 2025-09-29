# COMP.SE.140.Ex1
Exercise 01 for Continuous Development and Deployment - DevOps

## Running the Project

```bash
git clone -b exercise1 https://github.com/AbdullahShafqat-OG/COMP.SE.140.Ex1.git
cd COMP.SE.140.Ex1
docker-compose up --build
# … wait for ~10s
curl localhost:8199/status
docker-compose down
```

## Teacher's Instructions for Cleaning Up

```bash
docker-compose down --volumes
rm -f ./vstorage
```

## Reports and Logs
- [Report (PDF)](./Report.pdf)
- [Docker Status Output](./docker-status.txt)
- [LLM Usage Report](./llm.txt)