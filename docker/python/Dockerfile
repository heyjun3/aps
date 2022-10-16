FROM python:3.10.7

WORKDIR /usr/src/app

RUN apt-get update -y && apt-get install -y gcc python3-dev

COPY requirement.txt ./
RUN pip install --upgrade pip && \
pip install --no-cache -r requirement.txt

COPY . .
