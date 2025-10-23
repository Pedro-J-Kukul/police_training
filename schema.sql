--
-- PostgreSQL database dump
--

\restrict aP3abiEIxe15KhWWv0SZBxRGCbuDn8LUiIfPkz49yGK4eiTNh4lFbDG26hcTINm

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
    code text NOT NULL
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
-- Name: ranks; Type: TABLE; Schema: public; Owner: police
--

CREATE TABLE public.ranks (
    id bigint NOT NULL,
    rank text NOT NULL,
    code text NOT NULL
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
-- Name: users id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


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
-- Name: postings postings_code_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.postings
    ADD CONSTRAINT postings_code_key UNIQUE (code);


--
-- Name: postings postings_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.postings
    ADD CONSTRAINT postings_pkey PRIMARY KEY (id);


--
-- Name: postings postings_posting_key; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.postings
    ADD CONSTRAINT postings_posting_key UNIQUE (posting);


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
-- Name: idx_officers_user_id; Type: INDEX; Schema: public; Owner: police
--

CREATE INDEX idx_officers_user_id ON public.officers USING btree (user_id);


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

\unrestrict aP3abiEIxe15KhWWv0SZBxRGCbuDn8LUiIfPkz49yGK4eiTNh4lFbDG26hcTINm

