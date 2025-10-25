--
-- PostgreSQL database dump
--

\restrict sxlSOi4UVz6RJo3FuwEjPuu1r3tYck8pQJavjjHLCJbaWtJQGHOLB0FpxQPJLO5

-- Dumped from database version 17.6
-- Dumped by pg_dump version 17.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: citext; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: attendance_statuses; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.attendance_statuses (
    id integer NOT NULL,
    status text NOT NULL,
    counts_as_present boolean DEFAULT false NOT NULL
);


ALTER TABLE public.attendance_statuses OWNER TO police;

--
-- Name: attendance_statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.attendance_statuses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.attendance_statuses_id_seq OWNER TO police;

--
-- Name: attendance_statuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.attendance_statuses_id_seq OWNED BY public.attendance_statuses.id;


--
-- Name: enrollment_statuses; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.enrollment_statuses (
    id bigint NOT NULL,
    status text NOT NULL
);


ALTER TABLE public.enrollment_statuses OWNER TO police;

--
-- Name: enrollment_statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.enrollment_statuses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.enrollment_statuses_id_seq OWNER TO police;

--
-- Name: enrollment_statuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.enrollment_statuses_id_seq OWNED BY public.enrollment_statuses.id;


--
-- Name: formations; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.formations (
    id bigint NOT NULL,
    formation text NOT NULL,
    region_id bigint NOT NULL
);


ALTER TABLE public.formations OWNER TO police;

--
-- Name: formations_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.formations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.formations_id_seq OWNER TO police;

--
-- Name: formations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.formations_id_seq OWNED BY public.formations.id;


--
-- Name: officers; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.officers (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    regulation_number character varying(50) NOT NULL,
    rank_id bigint NOT NULL,
    posting_id bigint NOT NULL,
    formation_id bigint NOT NULL,
    region_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.officers OWNER TO police;

--
-- Name: officers_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

ALTER TABLE public.officers ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.officers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.permissions (
    id bigint NOT NULL,
    code text NOT NULL
);


ALTER TABLE public.permissions OWNER TO police;

--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.permissions_id_seq OWNER TO police;

--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: postings; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.postings (
    id bigint NOT NULL,
    posting text NOT NULL,
    code text
);


ALTER TABLE public.postings OWNER TO police;

--
-- Name: postings_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.postings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.postings_id_seq OWNER TO police;

--
-- Name: postings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.postings_id_seq OWNED BY public.postings.id;


--
-- Name: progress_statuses; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.progress_statuses (
    id integer NOT NULL,
    status text NOT NULL
);


ALTER TABLE public.progress_statuses OWNER TO police;

--
-- Name: progress_statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.progress_statuses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.progress_statuses_id_seq OWNER TO police;

--
-- Name: progress_statuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.progress_statuses_id_seq OWNED BY public.progress_statuses.id;


--
-- Name: ranks; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.ranks (
    id bigint NOT NULL,
    rank text NOT NULL,
    code text NOT NULL,
    annual_training_hours integer NOT NULL
);


ALTER TABLE public.ranks OWNER TO police;

--
-- Name: ranks_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.ranks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ranks_id_seq OWNER TO police;

--
-- Name: ranks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.ranks_id_seq OWNED BY public.ranks.id;


--
-- Name: regions; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.regions (
    id bigint NOT NULL,
    region text NOT NULL
);


ALTER TABLE public.regions OWNER TO police;

--
-- Name: regions_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.regions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.regions_id_seq OWNER TO police;

--
-- Name: regions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.regions_id_seq OWNED BY public.regions.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.roles (
    id bigint NOT NULL,
    role character varying(100) NOT NULL
);


ALTER TABLE public.roles OWNER TO police;

--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_seq OWNER TO police;

--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: roles_permissions; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.roles_permissions (
    permission_id bigint NOT NULL,
    role_id bigint NOT NULL
);


ALTER TABLE public.roles_permissions OWNER TO police;

--
-- Name: roles_users; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.roles_users (
    role_id bigint NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE public.roles_users OWNER TO police;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO police;

--
-- Name: tokens; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.tokens (
    hash bytea NOT NULL,
    user_id bigint NOT NULL,
    expiry timestamp without time zone NOT NULL,
    scope text NOT NULL
);


ALTER TABLE public.tokens OWNER TO police;

--
-- Name: training_categories; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.training_categories (
    id bigint NOT NULL,
    name text NOT NULL,
    is_active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.training_categories OWNER TO police;

--
-- Name: training_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.training_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.training_categories_id_seq OWNER TO police;

--
-- Name: training_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.training_categories_id_seq OWNED BY public.training_categories.id;


--
-- Name: training_enrollments; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.training_enrollments (
    id bigint NOT NULL,
    officer_id bigint NOT NULL,
    session_id bigint NOT NULL,
    enrollment_status_id bigint NOT NULL,
    attendance_status_id bigint,
    progress_status_id bigint NOT NULL,
    completion_date date,
    certificate_issued boolean DEFAULT false,
    certificate_number text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.training_enrollments OWNER TO police;

--
-- Name: training_enrollments_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.training_enrollments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.training_enrollments_id_seq OWNER TO police;

--
-- Name: training_enrollments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.training_enrollments_id_seq OWNED BY public.training_enrollments.id;


--
-- Name: training_sessions; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.training_sessions (
    id bigint NOT NULL,
    facilitator_id bigint NOT NULL,
    workshop_id bigint NOT NULL,
    formation_id bigint NOT NULL,
    region_id bigint NOT NULL,
    training_status_id bigint NOT NULL,
    session_date date NOT NULL,
    start_time time without time zone NOT NULL,
    end_time time without time zone NOT NULL,
    location text,
    max_capacity integer,
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.training_sessions OWNER TO police;

--
-- Name: training_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.training_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.training_sessions_id_seq OWNER TO police;

--
-- Name: training_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.training_sessions_id_seq OWNED BY public.training_sessions.id;


--
-- Name: training_status; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.training_status (
    id bigint NOT NULL,
    status text NOT NULL
);


ALTER TABLE public.training_status OWNER TO police;

--
-- Name: training_status_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.training_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.training_status_id_seq OWNER TO police;

--
-- Name: training_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.training_status_id_seq OWNED BY public.training_status.id;


--
-- Name: training_types; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.training_types (
    id bigint NOT NULL,
    name text NOT NULL,
    is_active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.training_types OWNER TO police;

--
-- Name: training_types_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.training_types_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.training_types_id_seq OWNER TO police;

--
-- Name: training_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.training_types_id_seq OWNED BY public.training_types.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    gender character(1) NOT NULL,
    email text NOT NULL,
    password_hash bytea NOT NULL,
    is_facilitator boolean DEFAULT false NOT NULL,
    is_officer boolean DEFAULT false NOT NULL,
    is_activated boolean DEFAULT false NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO police;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO police;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: workshops; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.workshops (
    id bigint NOT NULL,
    workshop_name text NOT NULL,
    category_id bigint NOT NULL,
    type_id bigint NOT NULL,
    credit_hours integer NOT NULL,
    description text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.workshops OWNER TO police;

--
-- Name: workshops_id_seq; Type: SEQUENCE; Schema: public; Owner: police
--

CREATE SEQUENCE public.workshops_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.workshops_id_seq OWNER TO police;

--
-- Name: workshops_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: police
--

ALTER SEQUENCE public.workshops_id_seq OWNED BY public.workshops.id;


--
-- Name: attendance_statuses id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.attendance_statuses ALTER COLUMN id SET DEFAULT nextval('public.attendance_statuses_id_seq'::regclass);


--
-- Name: enrollment_statuses id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.enrollment_statuses ALTER COLUMN id SET DEFAULT nextval('public.enrollment_statuses_id_seq'::regclass);


--
-- Name: formations id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.formations ALTER COLUMN id SET DEFAULT nextval('public.formations_id_seq'::regclass);


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: postings id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.postings ALTER COLUMN id SET DEFAULT nextval('public.postings_id_seq'::regclass);


--
-- Name: progress_statuses id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.progress_statuses ALTER COLUMN id SET DEFAULT nextval('public.progress_statuses_id_seq'::regclass);


--
-- Name: ranks id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.ranks ALTER COLUMN id SET DEFAULT nextval('public.ranks_id_seq'::regclass);


--
-- Name: regions id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.regions ALTER COLUMN id SET DEFAULT nextval('public.regions_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: training_categories id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_categories ALTER COLUMN id SET DEFAULT nextval('public.training_categories_id_seq'::regclass);


--
-- Name: training_enrollments id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_enrollments ALTER COLUMN id SET DEFAULT nextval('public.training_enrollments_id_seq'::regclass);


--
-- Name: training_sessions id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_sessions ALTER COLUMN id SET DEFAULT nextval('public.training_sessions_id_seq'::regclass);


--
-- Name: training_status id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_status ALTER COLUMN id SET DEFAULT nextval('public.training_status_id_seq'::regclass);


--
-- Name: training_types id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_types ALTER COLUMN id SET DEFAULT nextval('public.training_types_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: workshops id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.workshops ALTER COLUMN id SET DEFAULT nextval('public.workshops_id_seq'::regclass);


--
-- Name: attendance_statuses attendance_statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.attendance_statuses
    ADD CONSTRAINT attendance_statuses_pkey PRIMARY KEY (id);


--
-- Name: attendance_statuses attendance_statuses_status_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.attendance_statuses
    ADD CONSTRAINT attendance_statuses_status_key UNIQUE (status);


--
-- Name: enrollment_statuses enrollment_statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.enrollment_statuses
    ADD CONSTRAINT enrollment_statuses_pkey PRIMARY KEY (id);


--
-- Name: enrollment_statuses enrollment_statuses_status_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.enrollment_statuses
    ADD CONSTRAINT enrollment_statuses_status_key UNIQUE (status);


--
-- Name: formations formations_formation_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.formations
    ADD CONSTRAINT formations_formation_key UNIQUE (formation);


--
-- Name: formations formations_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.formations
    ADD CONSTRAINT formations_pkey PRIMARY KEY (id);


--
-- Name: officers officers_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.officers
    ADD CONSTRAINT officers_pkey PRIMARY KEY (id);


--
-- Name: officers officers_regulation_number_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.officers
    ADD CONSTRAINT officers_regulation_number_key UNIQUE (regulation_number);


--
-- Name: officers officers_user_id_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.officers
    ADD CONSTRAINT officers_user_id_key UNIQUE (user_id);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: postings postings_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.postings
    ADD CONSTRAINT postings_pkey PRIMARY KEY (id);


--
-- Name: progress_statuses progress_statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.progress_statuses
    ADD CONSTRAINT progress_statuses_pkey PRIMARY KEY (id);


--
-- Name: progress_statuses progress_statuses_status_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.progress_statuses
    ADD CONSTRAINT progress_statuses_status_key UNIQUE (status);


--
-- Name: ranks ranks_code_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.ranks
    ADD CONSTRAINT ranks_code_key UNIQUE (code);


--
-- Name: ranks ranks_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.ranks
    ADD CONSTRAINT ranks_pkey PRIMARY KEY (id);


--
-- Name: ranks ranks_rank_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.ranks
    ADD CONSTRAINT ranks_rank_key UNIQUE (rank);


--
-- Name: regions regions_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.regions
    ADD CONSTRAINT regions_pkey PRIMARY KEY (id);


--
-- Name: regions regions_region_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.regions
    ADD CONSTRAINT regions_region_key UNIQUE (region);


--
-- Name: roles_permissions roles_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT roles_permissions_pkey PRIMARY KEY (permission_id, role_id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: roles roles_role_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_role_key UNIQUE (role);


--
-- Name: roles_users roles_users_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles_users
    ADD CONSTRAINT roles_users_pkey PRIMARY KEY (role_id, user_id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (hash);


--
-- Name: training_categories training_categories_name_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_categories
    ADD CONSTRAINT training_categories_name_key UNIQUE (name);


--
-- Name: training_categories training_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_categories
    ADD CONSTRAINT training_categories_pkey PRIMARY KEY (id);


--
-- Name: training_enrollments training_enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_enrollments
    ADD CONSTRAINT training_enrollments_pkey PRIMARY KEY (id);


--
-- Name: training_sessions training_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_sessions
    ADD CONSTRAINT training_sessions_pkey PRIMARY KEY (id);


--
-- Name: training_status training_status_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_status
    ADD CONSTRAINT training_status_pkey PRIMARY KEY (id);


--
-- Name: training_status training_status_status_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_status
    ADD CONSTRAINT training_status_status_key UNIQUE (status);


--
-- Name: training_types training_types_name_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_types
    ADD CONSTRAINT training_types_name_key UNIQUE (name);


--
-- Name: training_types training_types_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_types
    ADD CONSTRAINT training_types_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: workshops workshops_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.workshops
    ADD CONSTRAINT workshops_pkey PRIMARY KEY (id);


--
-- Name: idx_officers_user_id; Type: INDEX; Schema: public; Owner: police
--

CREATE INDEX idx_officers_user_id ON public.officers USING btree (user_id);


--
-- Name: idx_training_enrollments_completion_date; Type: INDEX; Schema: public; Owner: police
--

CREATE INDEX idx_training_enrollments_completion_date ON public.training_enrollments USING btree (completion_date);


--
-- Name: idx_training_enrollments_officer_id; Type: INDEX; Schema: public; Owner: police
--

CREATE INDEX idx_training_enrollments_officer_id ON public.training_enrollments USING btree (officer_id);


--
-- Name: idx_training_enrollments_session_id; Type: INDEX; Schema: public; Owner: police
--

CREATE INDEX idx_training_enrollments_session_id ON public.training_enrollments USING btree (session_id);


--
-- Name: idx_training_enrollments_session_officer; Type: INDEX; Schema: public; Owner: police
--

CREATE UNIQUE INDEX idx_training_enrollments_session_officer ON public.training_enrollments USING btree (session_id, officer_id);


--
-- Name: workshops fk_category; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.workshops
    ADD CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES public.training_categories(id) ON DELETE RESTRICT;


--
-- Name: officers fk_formation; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.officers
    ADD CONSTRAINT fk_formation FOREIGN KEY (formation_id) REFERENCES public.formations(id) ON DELETE CASCADE;


--
-- Name: officers fk_posting; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.officers
    ADD CONSTRAINT fk_posting FOREIGN KEY (posting_id) REFERENCES public.postings(id) ON DELETE CASCADE;


--
-- Name: officers fk_rank; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.officers
    ADD CONSTRAINT fk_rank FOREIGN KEY (rank_id) REFERENCES public.ranks(id) ON DELETE CASCADE;


--
-- Name: formations fk_region; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.formations
    ADD CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES public.regions(id) ON DELETE CASCADE;


--
-- Name: officers fk_region; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.officers
    ADD CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES public.regions(id) ON DELETE CASCADE;


--
-- Name: training_enrollments fk_training_enrollments_attendance_status; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_enrollments
    ADD CONSTRAINT fk_training_enrollments_attendance_status FOREIGN KEY (attendance_status_id) REFERENCES public.attendance_statuses(id) ON DELETE CASCADE;


--
-- Name: training_enrollments fk_training_enrollments_enrollment_status; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_enrollments
    ADD CONSTRAINT fk_training_enrollments_enrollment_status FOREIGN KEY (enrollment_status_id) REFERENCES public.enrollment_statuses(id) ON DELETE CASCADE;


--
-- Name: training_enrollments fk_training_enrollments_officer; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_enrollments
    ADD CONSTRAINT fk_training_enrollments_officer FOREIGN KEY (officer_id) REFERENCES public.officers(id) ON DELETE CASCADE;


--
-- Name: training_enrollments fk_training_enrollments_progress_status; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_enrollments
    ADD CONSTRAINT fk_training_enrollments_progress_status FOREIGN KEY (progress_status_id) REFERENCES public.progress_statuses(id) ON DELETE CASCADE;


--
-- Name: training_enrollments fk_training_enrollments_session; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_enrollments
    ADD CONSTRAINT fk_training_enrollments_session FOREIGN KEY (session_id) REFERENCES public.training_sessions(id) ON DELETE CASCADE;


--
-- Name: workshops fk_type; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.workshops
    ADD CONSTRAINT fk_type FOREIGN KEY (type_id) REFERENCES public.training_types(id) ON DELETE RESTRICT;


--
-- Name: officers fk_user; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.officers
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: roles_permissions permission_roles_permissions; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT permission_roles_permissions FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: roles_permissions role_roles_permissions; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT role_roles_permissions FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: roles_users role_roles_users; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles_users
    ADD CONSTRAINT role_roles_users FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: training_sessions training_sessions_facilitator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_sessions
    ADD CONSTRAINT training_sessions_facilitator_id_fkey FOREIGN KEY (facilitator_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: training_sessions training_sessions_formation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_sessions
    ADD CONSTRAINT training_sessions_formation_id_fkey FOREIGN KEY (formation_id) REFERENCES public.formations(id) ON DELETE CASCADE;


--
-- Name: training_sessions training_sessions_region_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_sessions
    ADD CONSTRAINT training_sessions_region_id_fkey FOREIGN KEY (region_id) REFERENCES public.regions(id) ON DELETE CASCADE;


--
-- Name: training_sessions training_sessions_training_status_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_sessions
    ADD CONSTRAINT training_sessions_training_status_id_fkey FOREIGN KEY (training_status_id) REFERENCES public.training_status(id) ON DELETE CASCADE;


--
-- Name: training_sessions training_sessions_workshop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.training_sessions
    ADD CONSTRAINT training_sessions_workshop_id_fkey FOREIGN KEY (workshop_id) REFERENCES public.workshops(id) ON DELETE CASCADE;


--
-- Name: roles_users user_roles_users; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles_users
    ADD CONSTRAINT user_roles_users FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: tokens user_tokens; Type: FK CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT user_tokens FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict sxlSOi4UVz6RJo3FuwEjPuu1r3tYck8pQJavjjHLCJbaWtJQGHOLB0FpxQPJLO5

