### Database configuration

    CREATE USER go_domain_model WITH LOGIN CREATEDB PASSWORD 'password';
    CREATE DATABASE go_domain_model;
    ALTER DATABASE go_domain_model OWNER TO go_domain_model;


### Environment Variables

    EXPORT MODEL_DATABASE_DRIVER=postgres
    EXPORT MODEL_DATABASE_PARAMETERS=host=127.0.0.1 user=go_domain_model dbname=go_domain_model sslmode=require password=password