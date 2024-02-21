# Simple Banking System API

## Description

This is a simple banking system API that allows users to create accounts, deposit and withdraw money, and transfer money between accounts.

## Installation

1. Setting up postgres:

we are going to be using docker to set up the database. If you don't have docker installed, you can download it from [here](https://www.docker.com/products/docker-desktop)

```bash
Make postgres
```

2. Creating the databse

Now that we have a postgres instance, we need to create the db. We can do this by running the following command:

```bash
Make createdb
```

3. Create the database schema

Finally we need to create the schema for the database. For that we can run the following command:

```bash
Make migrateup
```
