-- Users queries
-- name: GetUserByID :one
SELECT id, email, name, password_hash, google_id, access_token, refresh_token, auth_type, is_active, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByGoogleID :one
SELECT id, email, name, password_hash, google_id, access_token, refresh_token, auth_type, is_active, created_at, updated_at
FROM users
WHERE google_id = $1;

-- name: GetUserByEmail :one
SELECT id, email, name, password_hash, google_id, access_token, refresh_token, auth_type, is_active, created_at, updated_at
FROM users
WHERE email = $1;

-- name: CreateUserWithPassword :one
INSERT INTO users (email, name, password_hash, auth_type)
VALUES ($1, $2, $3, 'password')
RETURNING *;

-- name: CreateUserWithGoogle :one
INSERT INTO users (email, name, google_id, access_token, refresh_token, auth_type)
VALUES ($1, $2, $3, $4, $5, 'google')
RETURNING *;

-- name: CreateUserWithBoth :one
INSERT INTO users (email, name, password_hash, google_id, access_token, refresh_token, auth_type)
VALUES ($1, $2, $3, $4, $5, $6, 'both')
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users 
SET password_hash = $2, auth_type = CASE 
  WHEN auth_type = 'google' THEN 'both' 
  ELSE 'password' 
END, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserGoogleAuth :exec
UPDATE users 
SET google_id = $2, access_token = $3, refresh_token = $4, auth_type = CASE 
  WHEN auth_type = 'password' THEN 'both' 
  ELSE 'google' 
END, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserTokens :exec
UPDATE users 
SET access_token = $2, refresh_token = $3, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users 
SET name = $2, email = $3, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserActiveStatus :exec
UPDATE users 
SET is_active = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, name, auth_type, is_active, created_at, updated_at
FROM users
ORDER BY created_at DESC;

-- Companies queries
-- name: GetCompanyByID :one
SELECT id, name, linkedin_url, industry, created_at, updated_at
FROM companies
WHERE id = $1;

-- name: GetCompanyByName :one
SELECT id, name, linkedin_url, industry, created_at, updated_at
FROM companies
WHERE name = $1;

-- name: GetCompanyByLinkedInURL :one
SELECT id, name, linkedin_url, industry, created_at, updated_at
FROM companies
WHERE linkedin_url = $1;

-- name: CreateCompany :one
INSERT INTO companies (name, linkedin_url, industry)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCompany :exec
UPDATE companies 
SET name = $2, linkedin_url = $3, industry = $4, updated_at = NOW()
WHERE id = $1;

-- name: ListCompanies :many
SELECT id, name, linkedin_url, industry, created_at, updated_at
FROM companies
ORDER BY name;

-- LinkedIn Profiles queries
-- name: GetLinkedInProfileByID :one
SELECT id, linkedin_url, name, location, current_company_id, headline, created_at, updated_at
FROM linkedin_profiles
WHERE id = $1;

-- name: GetLinkedInProfileByURL :one
SELECT id, linkedin_url, name, location, current_company_id, headline, created_at, updated_at
FROM linkedin_profiles
WHERE linkedin_url = $1;

-- name: CreateLinkedInProfile :one
INSERT INTO linkedin_profiles (linkedin_url, name, location, current_company_id, headline)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateLinkedInProfile :exec
UPDATE linkedin_profiles 
SET name = $2, location = $3, current_company_id = $4, headline = $5, updated_at = NOW()
WHERE id = $1;

-- name: ListLinkedInProfiles :many
SELECT id, linkedin_url, name, location, current_company_id, headline, created_at, updated_at
FROM linkedin_profiles
ORDER BY name;

-- Profile Companies (Employment History) queries
-- name: GetProfileCompanies :many
SELECT pc.id, pc.profile_id, pc.company_id, pc.position, pc.start_date, pc.end_date, pc.is_current, pc.created_at,
       c.name as company_name, c.linkedin_url as company_linkedin_url
FROM profile_companies pc
JOIN companies c ON pc.company_id = c.id
WHERE pc.profile_id = $1
ORDER BY pc.start_date DESC;

-- name: CreateProfileCompany :one
INSERT INTO profile_companies (profile_id, company_id, position, start_date, end_date, is_current)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProfileCompany :exec
UPDATE profile_companies 
SET position = $3, start_date = $4, end_date = $5, is_current = $6
WHERE profile_id = $1 AND company_id = $2;

-- name: DeleteProfileCompany :exec
DELETE FROM profile_companies 
WHERE profile_id = $1 AND company_id = $2;

-- Tracked Connections queries
-- name: GetTrackedConnections :many
SELECT tc.id, tc.user_id, tc.profile_id, tc.created_at, tc.last_checked_at,
       lp.linkedin_url, lp.name, lp.location, lp.headline,
       c.name as company_name
FROM tracked_connections tc
JOIN linkedin_profiles lp ON tc.profile_id = lp.id
LEFT JOIN companies c ON lp.current_company_id = c.id
WHERE tc.user_id = $1
ORDER BY tc.last_checked_at DESC;

-- name: GetTrackedConnection :one
SELECT tc.id, tc.user_id, tc.profile_id, tc.created_at, tc.last_checked_at
FROM tracked_connections tc
WHERE tc.user_id = $1 AND tc.profile_id = $2;

-- name: CreateTrackedConnection :one
INSERT INTO tracked_connections (user_id, profile_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateTrackedConnectionLastChecked :exec
UPDATE tracked_connections 
SET last_checked_at = NOW()
WHERE id = $1;

-- name: DeleteTrackedConnection :exec
DELETE FROM tracked_connections 
WHERE user_id = $1 AND profile_id = $2;

-- Connection Relationships queries
-- name: GetConnectionRelationships :many
SELECT cr.id, cr.profile_a_id, cr.profile_b_id, cr.degree, cr.discovered_at, cr.discovered_by_user_id,
       lp1.name as profile_a_name, lp1.linkedin_url as profile_a_url,
       lp2.name as profile_b_name, lp2.linkedin_url as profile_b_url
FROM connection_relationships cr
JOIN linkedin_profiles lp1 ON cr.profile_a_id = lp1.id
JOIN linkedin_profiles lp2 ON cr.profile_b_id = lp2.id
WHERE cr.profile_a_id = $1 OR cr.profile_b_id = $1
ORDER BY cr.degree, cr.discovered_at DESC;

-- name: GetConnectionsByDegree :many
SELECT cr.id, cr.profile_a_id, cr.profile_b_id, cr.degree, cr.discovered_at, cr.discovered_by_user_id,
       lp1.name as profile_a_name, lp1.linkedin_url as profile_a_url,
       lp2.name as profile_b_name, lp2.linkedin_url as profile_b_url
FROM connection_relationships cr
JOIN linkedin_profiles lp1 ON cr.profile_a_id = lp1.id
JOIN linkedin_profiles lp2 ON cr.profile_b_id = lp2.id
WHERE (cr.profile_a_id = $1 OR cr.profile_b_id = $1) AND cr.degree = $2
ORDER BY cr.discovered_at DESC;

-- name: CreateConnectionRelationship :one
INSERT INTO connection_relationships (profile_a_id, profile_b_id, degree, discovered_by_user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CheckConnectionExists :one
SELECT id FROM connection_relationships 
WHERE ((profile_a_id = $1 AND profile_b_id = $2) OR (profile_a_id = $2 AND profile_b_id = $1)) 
AND degree = $3;

-- Automation Rules queries
-- name: GetAutomationRules :many
SELECT id, user_id, name, company_filter, location_filter, action_type, message_template, is_active, created_at
FROM automation_rules
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetAutomationRuleByID :one
SELECT id, user_id, name, company_filter, location_filter, action_type, message_template, is_active, created_at
FROM automation_rules
WHERE id = $1 AND user_id = $2;

-- name: CreateAutomationRule :one
INSERT INTO automation_rules (user_id, name, company_filter, location_filter, action_type, message_template)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateAutomationRule :exec
UPDATE automation_rules 
SET name = $3, company_filter = $4, location_filter = $5, action_type = $6, message_template = $7, is_active = $8
WHERE id = $1 AND user_id = $2;

-- name: DeleteAutomationRule :exec
DELETE FROM automation_rules 
WHERE id = $1 AND user_id = $2;

-- name: GetActiveAutomationRules :many
SELECT id, user_id, name, company_filter, location_filter, action_type, message_template, is_active, created_at
FROM automation_rules
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- Business Logic queries
-- name: GetNewConnectionsForUser :many
SELECT DISTINCT lp.id, lp.linkedin_url, lp.name, lp.location, lp.headline, lp.created_at,
       c.name as company_name,
       MIN(cr.degree) as connection_degree
FROM linkedin_profiles lp
LEFT JOIN companies c ON lp.current_company_id = c.id
JOIN connection_relationships cr ON (cr.profile_a_id = lp.id OR cr.profile_b_id = lp.id)
JOIN tracked_connections tc ON (
    (cr.profile_a_id = tc.profile_id AND cr.profile_b_id = lp.id) OR
    (cr.profile_b_id = tc.profile_id AND cr.profile_a_id = lp.id)
)
WHERE tc.user_id = $1 
  AND lp.id NOT IN (
    SELECT DISTINCT tc2.profile_id 
    FROM tracked_connections tc2 
    WHERE tc2.user_id = $1
  )
  AND cr.discovered_at > $2
GROUP BY lp.id, lp.linkedin_url, lp.name, lp.location, lp.headline, lp.created_at, c.name
ORDER BY lp.created_at DESC;

-- name: GetProfilesMatchingRules :many
SELECT lp.id, lp.linkedin_url, lp.name, lp.location, lp.headline,
       c.name as company_name,
       ar.id as rule_id, ar.name as rule_name, ar.action_type, ar.message_template
FROM linkedin_profiles lp
LEFT JOIN companies c ON lp.current_company_id = c.id
CROSS JOIN automation_rules ar
WHERE ar.user_id = $1 
  AND ar.is_active = true
  AND (ar.company_filter IS NULL OR c.name ILIKE '%' || ar.company_filter || '%')
  AND (ar.location_filter IS NULL OR lp.location ILIKE '%' || ar.location_filter || '%')
  AND lp.id NOT IN (
    SELECT DISTINCT tc.profile_id 
    FROM tracked_connections tc 
    WHERE tc.user_id = $1
  );

-- Utility queries
-- name: PingDb :one
SELECT 1 as result;
