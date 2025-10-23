--
-- PostgreSQL database dump
--

\restrict EQo3Y6erx3py6rcI8hifm4EeSAd9c6e46b4bX8q0BaexfANavIfohe2CykOu2ws

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
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: police
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


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

\unrestrict EQo3Y6erx3py6rcI8hifm4EeSAd9c6e46b4bX8q0BaexfANavIfohe2CykOu2ws

