--
-- PostgreSQL database dump
--

-- Dumped from database version 11.18
-- Dumped by pg_dump version 11.18

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: access_level; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.access_level (
    application character varying,
    feature character varying,
    accesslevel character varying
);


ALTER TABLE public.access_level OWNER TO postgres;

--
-- Data for Name: access_level; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.access_level (application, feature, accesslevel) FROM stdin;
User Management	Login	2
User Management	User registration	2
User Management	User registration	0
Restsim Offerings	Product Documentation	1
Restsim Offerings	User Documentation	2
Restsim Offerings	Simulation Catalog	1
Restsim Offerings	Self Service	1
Insights	Deployments Status summary	2
Insights	Specific Deployment status	1
Ticket Management	New Sim request	1
Ticket Management	Support Request	1
Ticket Management	Information Request	1
Communication	News Feed	2
Communication	Announcement	2
Communication	Message Broadcast	0
Search	Site Navigation	2
\.


--
-- PostgreSQL database dump complete
--

