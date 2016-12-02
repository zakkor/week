--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.1
-- Dumped by pg_dump version 9.6.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: comments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE comments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    parent_post uuid,
    author character varying(64) NOT NULL,
    date timestamp without time zone NOT NULL,
    content character varying(3000)
);


--
-- Name: posts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE posts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(40) NOT NULL,
    author character varying(64) NOT NULL,
    date timestamp without time zone NOT NULL,
    content character varying(3000)
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE users (
    email character varying(64) NOT NULL,
    name character varying(64) NOT NULL,
    password character varying(64) NOT NULL,
    session character varying(64)
);


--
-- Data for Name: comments; Type: TABLE DATA; Schema: public; Owner: -
--

COPY comments (id, parent_post, author, date, content) FROM stdin;
e0d9edb3-b986-4470-9304-36f3d076885e	7e39c106-d80e-4e99-8a6e-74a40974c25a	edy	2016-12-01 23:29:48.6023	asdasdasd
c788f11e-2456-41a3-b62c-5b4feaca4a2b	7e39c106-d80e-4e99-8a6e-74a40974c25a	edy	2016-12-01 23:30:00.683351	fgdfgdfg
12846e4b-a584-40c8-b163-2c280c89946e	7e39c106-d80e-4e99-8a6e-74a40974c25a	edy	2016-12-01 23:30:30.662743	aasdasd
e8a7912d-ac40-4127-a319-c6964bc7efc4	7e39c106-d80e-4e99-8a6e-74a40974c25a	edy	2016-12-02 03:17:54.848013	a
\.


--
-- Data for Name: posts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY posts (id, title, author, date, content) FROM stdin;
7e39c106-d80e-4e99-8a6e-74a40974c25a	asdasdasdas	edy	2016-12-01 23:29:36.444002	asadsasd
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY users (email, name, password, session) FROM stdin;
asd@asd	edy	31d3dd0d209225fe57cf87e70bc562074a45037bc3652ea50f7d8746014c10ef	c594a0f90bd6291cd1f70e70b69fb27b908e701a567ad51e2d48cd4193612971
\.


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: comments; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE comments TO ed;


--
-- Name: posts; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE posts TO ed;


--
-- Name: users; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE users TO ed;


--
-- PostgreSQL database dump complete
--

