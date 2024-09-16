CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);
CREATE TYPE tender_status AS ENUM(
    'CREATED',
    'PUBLISHED',
    'CLOSED'
);
CREATE TABLE tender(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    name VARCHAR(100),
    description TEXT,
    service_type VARCHAR(100),
    status tender_status DEFAULT 'CREATED',
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,   
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TYPE authortype AS ENUM(
    'User',
    'Organization'      
);
CREATE TYPE bid_status AS ENUM(
    'CREATED',
    'PUBLISHED',
    'CANCELED'  
);
CREATE TABLE bid(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100),
    description TEXT,
    status bid_status DEFAULT 'CREATED',
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_type authortype, 
    author_id UUID,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE tender_history (
    id UUID PRIMARY KEY,
    tender_id UUID NOT NULL,
    name TEXT,
    description TEXT,
    service_type TEXT,
    version INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tender_id) REFERENCES tender(id)
);

CREATE TABLE bid_history (
    id UUID PRIMARY KEY,
    bid_id UUID NOT NULL,
    name TEXT,
    description TEXT,
    version INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bid_id) REFERENCES bid(id)
);
CREATE OR REPLACE FUNCTION update_tender_version()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO tender_history (id, tender_id, name, description, service_type, version)
    VALUES (uuid_generate_v4(), OLD.id, OLD.name, OLD.description, OLD.service_type, OLD.version);

    IF NEW.name <> OLD.name OR 
       NEW.description <> OLD.description OR 
       NEW.service_type <> OLD.service_type THEN
        NEW.version := OLD.version + 1;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tender_version_trigger
BEFORE UPDATE ON tender
FOR EACH ROW EXECUTE FUNCTION update_tender_version();

CREATE OR REPLACE FUNCTION update_bid_version()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO bid_history (id, bid_id, name, description, version)
    VALUES (uuid_generate_v4(), OLD.id, OLD.name, OLD.description, OLD.version);
    IF NEW.name <> OLD.name OR 
       NEW.description <> OLD.description THEN
        NEW.version := OLD.version + 1;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER bid_version_trigger
BEFORE UPDATE ON bid
FOR EACH ROW EXECUTE FUNCTION update_bid_version();


