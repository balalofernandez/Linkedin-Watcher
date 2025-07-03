-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (supports both OAuth2 and traditional authentication)
CREATE TABLE users (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email             VARCHAR(255) UNIQUE NOT NULL,
  name              VARCHAR(255) NOT NULL,
  password_hash     VARCHAR(255), -- For traditional username/password auth
  google_id         VARCHAR(255) UNIQUE, -- For OAuth2 (optional)
  access_token      TEXT, -- For OAuth2 (optional)
  refresh_token     TEXT, -- For OAuth2 (optional)
  auth_type         VARCHAR(20) NOT NULL DEFAULT 'password' CHECK (auth_type IN ('password', 'google', 'both')),
  is_active         BOOLEAN NOT NULL DEFAULT TRUE,
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT valid_auth_type CHECK (
    (auth_type = 'password' AND password_hash IS NOT NULL) OR
    (auth_type = 'google' AND google_id IS NOT NULL) OR
    (auth_type = 'both' AND password_hash IS NOT NULL AND google_id IS NOT NULL)
  )
);

-- Companies table
CREATE TABLE companies (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name              VARCHAR(255) UNIQUE NOT NULL,
  linkedin_url      VARCHAR(500) UNIQUE,
  industry          VARCHAR(255),
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

-- LinkedIn profiles table
CREATE TABLE linkedin_profiles (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  linkedin_url      VARCHAR(500) UNIQUE NOT NULL,
  linkedin_id       VARCHAR(255),
  name              VARCHAR(255) NOT NULL,
  location          VARCHAR(255),
  current_company_id UUID REFERENCES companies(id),
  headline          VARCHAR(500),
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT valid_linkedin_url CHECK (linkedin_url ~ '^https?://(www\.)?linkedin\.com/in/[a-zA-Z0-9-]+/?$')
);

-- Employment history table
CREATE TABLE profile_companies (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  profile_id        UUID NOT NULL REFERENCES linkedin_profiles(id) ON DELETE CASCADE,
  company_id        UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
  position          VARCHAR(255) NOT NULL,
  start_date        DATE NOT NULL,
  end_date          DATE,
  is_current        BOOLEAN NOT NULL DEFAULT FALSE,
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(profile_id, company_id, position, start_date)
);

-- Tracked connections table (what users are monitoring)
CREATE TABLE tracked_connections (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  profile_id        UUID NOT NULL REFERENCES linkedin_profiles(id) ON DELETE CASCADE,
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  last_checked_at   TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, profile_id)
);

-- Connection relationships table (graph structure with degree)
CREATE TABLE connection_relationships (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  profile_a_id      UUID NOT NULL REFERENCES linkedin_profiles(id) ON DELETE CASCADE,
  profile_b_id      UUID NOT NULL REFERENCES linkedin_profiles(id) ON DELETE CASCADE,
  degree            INTEGER NOT NULL CHECK (degree IN (1, 2, 3)),
  discovered_at     TIMESTAMP NOT NULL DEFAULT NOW(),
  discovered_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(profile_a_id, profile_b_id, degree)
);

-- Automation rules table
CREATE TABLE automation_rules (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name              VARCHAR(255) NOT NULL,
  company_filter    VARCHAR(255),
  location_filter   VARCHAR(255),
  action_type       VARCHAR(50) NOT NULL CHECK (action_type IN ('send_connection_request', 'save_profile', 'notify')),
  message_template  TEXT,
  is_active         BOOLEAN NOT NULL DEFAULT TRUE,
  created_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_linkedin_profiles_url ON linkedin_profiles(linkedin_url);
CREATE INDEX idx_linkedin_profiles_company ON linkedin_profiles(current_company_id);
CREATE INDEX idx_profile_companies_profile ON profile_companies(profile_id);
CREATE INDEX idx_profile_companies_company ON profile_companies(company_id);
CREATE INDEX idx_tracked_connections_user ON tracked_connections(user_id);
CREATE INDEX idx_tracked_connections_profile ON tracked_connections(profile_id);
CREATE INDEX idx_connection_relationships_profile_a ON connection_relationships(profile_a_id);
CREATE INDEX idx_connection_relationships_profile_b ON connection_relationships(profile_b_id);
CREATE INDEX idx_connection_relationships_degree ON connection_relationships(degree);
CREATE INDEX idx_automation_rules_user ON automation_rules(user_id);
CREATE INDEX idx_automation_rules_active ON automation_rules(is_active);
