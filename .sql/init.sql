--
-- TOC entry 443 (class 1255 OID 17233)
-- Name: features(); Type: FUNCTION; Schema: public; Owner: supabase_admin
--

CREATE FUNCTION public.features() RETURNS json
    LANGUAGE sql
    AS $$select json_agg(distinct ( 
  case
    when (r.feature != '') then r.feature
    else r.parent_suite
  end)) as features
from results as r
  inner join launches as l on r.launch_id = l.id
  inner join versions as v on l.version_id = v.id
  inner join projects as p on v.project_id = p.id
where p.id = 1 and l.is_template = true$$;


ALTER FUNCTION public.features() OWNER TO supabase_admin;

--
-- TOC entry 444 (class 1255 OID 17234)
-- Name: results(bigint); Type: FUNCTION; Schema: public; Owner: supabase_admin
--

CREATE FUNCTION public.results(version bigint) RETURNS json
    LANGUAGE sql
    AS $$select json_agg(j) from (select r.name as name, r.feature, r.parent_suite, r.suite, r.sub_suite, r.status, 
  l.name as launch_name, l.id as launch_id, l.created_at
from results as r
  join launches as l on r.launch_id = l.id
where l.id in (
  select ll.id 
  from launches as ll
  where ll.version_id = (version)
  order by ll.created_at desc
  limit 1
  )
) as j$$;


ALTER FUNCTION public.results(version bigint) OWNER TO supabase_admin;

--
-- TOC entry 442 (class 1255 OID 17227)
-- Name: template(bigint); Type: FUNCTION; Schema: public; Owner: supabase_admin
--

CREATE FUNCTION public.template(version bigint) RETURNS json
    LANGUAGE sql
    AS $$select json_agg(r.*)
from results as r
  inner join launches as l on r.launch_id = l.id 
  inner join versions as v on l.version_id = v.id 
  inner join projects as p on v.project_id = p.id 
where l.is_template = true and p.id in (select versions.project_id from versions where versions.id = (version))$$;


ALTER FUNCTION public.template(version bigint) OWNER TO supabase_admin;

--
-- TOC entry 445 (class 1255 OID 17231)
-- Name: version(); Type: FUNCTION; Schema: public; Owner: supabase_admin
--

CREATE FUNCTION public.version() RETURNS json
    LANGUAGE sql
    AS $$select json_agg(v) from (
select distinct on (v.id) 
  v.id, v.version_name, v.repo, count(distinct l.id) as launch_count, v.created_at,
  sq.last_launch_id, sq.last_launch_at, sq.is_template_launch, sq.launch_name, sq.results_number, 
  sq.passed_number, sq.skipped_number, sq.failed_number, sq.undefined_number
from versions as v
  left join launches as l on l.version_id = v.id
  left join (
    select ll.id as last_launch_id, ll.created_at as last_launch_at, vv.id as vv_id, 
    ll.is_template as is_template_launch, ll.name as launch_name,
    count(distinct r.id) as results_number, rp.cnt as passed_number, rs.cnt as skipped_number,
    rf.cnt as failed_number, ru.cnt as undefined_number
    from launches as ll 
      left join versions as vv on vv.id = ll.version_id
      left join results as r on r.launch_id = ll.id
      left join (
        select count(distinct rr.id) as cnt, rr.launch_id 
        from results as rr where rr.status = 'passed' 
        group by rr.launch_id
      ) as rp on rp.launch_id = ll.id
      left join (
        select count(distinct rr.id) as cnt, rr.launch_id 
        from results as rr where rr.status = 'skipped' 
        group by rr.launch_id
      ) as rs on rs.launch_id = ll.id
      left join (
        select count(distinct rr.id) as cnt, rr.launch_id 
        from results as rr where rr.status in ('failed', 'error')
        group by rr.launch_id
      ) as rf on rf.launch_id = ll.id
      left join (
        select count(distinct rr.id) as cnt, rr.launch_id 
        from results as rr where rr.status not in ('passed', 'skipped', 'failed', 'error')
        group by rr.launch_id
      ) as ru on ru.launch_id = ll.id
    group by ll.id, vv.id, rp.cnt, rs.cnt, rf.cnt, ru.cnt
    order by ll.created_at desc
  ) as sq on sq.vv_id = v.id
where v.project_id = 1
group by v.id, sq.last_launch_id, sq.last_launch_at, sq.is_template_launch, sq.launch_name, sq.results_number, 
  sq.passed_number, sq.skipped_number, sq.failed_number, sq.undefined_number
order by v.id, sq.last_launch_at desc) as v$$;


ALTER FUNCTION public.version() OWNER TO supabase_admin;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 253 (class 1259 OID 17148)
-- Name: versions; Type: TABLE; Schema: public; Owner: supabase_admin
--

CREATE TABLE public.versions (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    version_name character varying NOT NULL,
    repo character varying NOT NULL,
    project_id bigint NOT NULL,
    secondary_id character varying DEFAULT 'suite'::character varying
);


ALTER TABLE public.versions OWNER TO supabase_admin;

--
-- TOC entry 2906 (class 0 OID 0)
-- Dependencies: 253
-- Name: COLUMN versions.secondary_id; Type: COMMENT; Schema: public; Owner: supabase_admin
--

COMMENT ON COLUMN public.versions.secondary_id IS 'what to use as second identifier (first - name) may be ''suite'', ''parent_suite'' or ''sub_suite''';


--
-- TOC entry 254 (class 1259 OID 17151)
-- Name: clients_id_seq; Type: SEQUENCE; Schema: public; Owner: supabase_admin
--

ALTER TABLE public.versions ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.clients_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 252 (class 1259 OID 17134)
-- Name: launches; Type: TABLE; Schema: public; Owner: supabase_admin
--

CREATE TABLE public.launches (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    is_template boolean DEFAULT false NOT NULL,
    origin character varying,
    user_id uuid NOT NULL,
    version_id bigint NOT NULL,
    duration integer,
    name character varying NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.launches OWNER TO supabase_admin;

--
-- TOC entry 258 (class 1259 OID 17218)
-- Name: launches_id_seq; Type: SEQUENCE; Schema: public; Owner: supabase_admin
--

ALTER TABLE public.launches ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.launches_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 255 (class 1259 OID 17165)
-- Name: projects; Type: TABLE; Schema: public; Owner: supabase_admin
--

CREATE TABLE public.projects (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    name character varying NOT NULL,
    description text,
    owner_id uuid
);


ALTER TABLE public.projects OWNER TO supabase_admin;

--
-- TOC entry 256 (class 1259 OID 17168)
-- Name: projects_id_seq; Type: SEQUENCE; Schema: public; Owner: supabase_admin
--

ALTER TABLE public.projects ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.projects_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 257 (class 1259 OID 17182)
-- Name: results; Type: TABLE; Schema: public; Owner: supabase_admin
--

CREATE TABLE public.results (
    id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    name character varying NOT NULL,
    suite character varying NOT NULL,
    feature character varying NOT NULL,
    status character varying,
    description text,
    steps json,
    duration integer,
    fullname character varying,
    launch_id bigint NOT NULL,
    parent_suite character varying,
    sub_suite character varying
);


ALTER TABLE public.results OWNER TO supabase_admin;

--
-- TOC entry 2725 (class 2606 OID 17159)
-- Name: versions clients_pkey; Type: CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.versions
    ADD CONSTRAINT clients_pkey PRIMARY KEY (id);


--
-- TOC entry 2723 (class 2606 OID 17211)
-- Name: launches launches_pkey; Type: CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.launches
    ADD CONSTRAINT launches_pkey PRIMARY KEY (id);


--
-- TOC entry 2727 (class 2606 OID 17176)
-- Name: projects projects_pkey; Type: CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT projects_pkey PRIMARY KEY (id);


--
-- TOC entry 2729 (class 2606 OID 17189)
-- Name: results results_pkey; Type: CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT results_pkey PRIMARY KEY (id);


--
-- TOC entry 2732 (class 2606 OID 17177)
-- Name: versions clients_project_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.versions
    ADD CONSTRAINT clients_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.projects(id);


--
-- TOC entry 2730 (class 2606 OID 17143)
-- Name: launches launches_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.launches
    ADD CONSTRAINT launches_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users(id);


--
-- TOC entry 2731 (class 2606 OID 17202)
-- Name: launches launches_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.launches
    ADD CONSTRAINT launches_version_id_fkey FOREIGN KEY (version_id) REFERENCES public.versions(id);


--
-- TOC entry 2733 (class 2606 OID 17235)
-- Name: projects projects_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT projects_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES auth.users(id);


--
-- TOC entry 2734 (class 2606 OID 17219)
-- Name: results results_launch_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: supabase_admin
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT results_launch_id_fkey FOREIGN KEY (launch_id) REFERENCES public.launches(id);


--
-- TOC entry 2888 (class 3256 OID 17196)
-- Name: launches Enable access to all users; Type: POLICY; Schema: public; Owner: supabase_admin
--

CREATE POLICY "Enable access to all users" ON public.launches FOR SELECT USING (true);


--
-- TOC entry 2889 (class 3256 OID 17197)
-- Name: projects Enable access to all users; Type: POLICY; Schema: public; Owner: supabase_admin
--

CREATE POLICY "Enable access to all users" ON public.projects FOR SELECT USING (true);


--
-- TOC entry 2890 (class 3256 OID 17198)
-- Name: results Enable access to all users; Type: POLICY; Schema: public; Owner: supabase_admin
--

CREATE POLICY "Enable access to all users" ON public.results FOR SELECT USING (true);


--
-- TOC entry 2887 (class 3256 OID 17195)
-- Name: versions Enable access to all users; Type: POLICY; Schema: public; Owner: supabase_admin
--

CREATE POLICY "Enable access to all users" ON public.versions FOR SELECT USING (true);


--
-- TOC entry 2891 (class 3256 OID 17201)
-- Name: launches Enable update for users by id; Type: POLICY; Schema: public; Owner: supabase_admin
--

CREATE POLICY "Enable update for users by id" ON public.launches FOR UPDATE USING ((auth.uid() = user_id)) WITH CHECK ((auth.uid() = user_id));


--
-- TOC entry 2893 (class 3256 OID 17199)
-- Name: launches insert allowed only by yourself; Type: POLICY; Schema: public; Owner: supabase_admin
--

CREATE POLICY "insert allowed only by yourself" ON public.launches FOR INSERT WITH CHECK (((auth.uid() = user_id) AND ((is_template = false) OR (auth.uid() IN ( SELECT p.owner_id
   FROM ((public.projects p
     JOIN public.versions v ON ((v.project_id = p.id)))
     JOIN public.launches l ON ((l.version_id = v.id))))))));


--
-- TOC entry 2892 (class 3256 OID 17200)
-- Name: results insert allowed only by yourself; Type: POLICY; Schema: public; Owner: supabase_admin
--

CREATE POLICY "insert allowed only by yourself" ON public.results FOR INSERT WITH CHECK ((auth.uid() IN ( SELECT launches.user_id
   FROM public.launches
  WHERE (launches.id = results.launch_id))));


--
-- TOC entry 2883 (class 0 OID 17134)
-- Dependencies: 252
-- Name: launches; Type: ROW SECURITY; Schema: public; Owner: supabase_admin
--

ALTER TABLE public.launches ENABLE ROW LEVEL SECURITY;

--
-- TOC entry 2885 (class 0 OID 17165)
-- Dependencies: 255
-- Name: projects; Type: ROW SECURITY; Schema: public; Owner: supabase_admin
--

ALTER TABLE public.projects ENABLE ROW LEVEL SECURITY;

--
-- TOC entry 2886 (class 0 OID 17182)
-- Dependencies: 257
-- Name: results; Type: ROW SECURITY; Schema: public; Owner: supabase_admin
--

ALTER TABLE public.results ENABLE ROW LEVEL SECURITY;

--
-- TOC entry 2884 (class 0 OID 17148)
-- Dependencies: 253
-- Name: versions; Type: ROW SECURITY; Schema: public; Owner: supabase_admin
--

ALTER TABLE public.versions ENABLE ROW LEVEL SECURITY;

--
-- TOC entry 2901 (class 0 OID 0)
-- Dependencies: 8
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT USAGE ON SCHEMA public TO anon;
GRANT USAGE ON SCHEMA public TO authenticated;
GRANT USAGE ON SCHEMA public TO service_role;


--
-- TOC entry 2902 (class 0 OID 0)
-- Dependencies: 443
-- Name: FUNCTION features(); Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON FUNCTION public.features() TO postgres;
GRANT ALL ON FUNCTION public.features() TO anon;
GRANT ALL ON FUNCTION public.features() TO authenticated;
GRANT ALL ON FUNCTION public.features() TO service_role;


--
-- TOC entry 2903 (class 0 OID 0)
-- Dependencies: 444
-- Name: FUNCTION results(version bigint); Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON FUNCTION public.results(version bigint) TO postgres;
GRANT ALL ON FUNCTION public.results(version bigint) TO anon;
GRANT ALL ON FUNCTION public.results(version bigint) TO authenticated;
GRANT ALL ON FUNCTION public.results(version bigint) TO service_role;


--
-- TOC entry 2904 (class 0 OID 0)
-- Dependencies: 442
-- Name: FUNCTION template(version bigint); Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON FUNCTION public.template(version bigint) TO postgres;
GRANT ALL ON FUNCTION public.template(version bigint) TO anon;
GRANT ALL ON FUNCTION public.template(version bigint) TO authenticated;
GRANT ALL ON FUNCTION public.template(version bigint) TO service_role;


--
-- TOC entry 2905 (class 0 OID 0)
-- Dependencies: 445
-- Name: FUNCTION version(); Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON FUNCTION public.version() TO postgres;
GRANT ALL ON FUNCTION public.version() TO anon;
GRANT ALL ON FUNCTION public.version() TO authenticated;
GRANT ALL ON FUNCTION public.version() TO service_role;


--
-- TOC entry 2907 (class 0 OID 0)
-- Dependencies: 253
-- Name: TABLE versions; Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON TABLE public.versions TO postgres;
GRANT ALL ON TABLE public.versions TO anon;
GRANT ALL ON TABLE public.versions TO authenticated;
GRANT ALL ON TABLE public.versions TO service_role;


--
-- TOC entry 2908 (class 0 OID 0)
-- Dependencies: 254
-- Name: SEQUENCE clients_id_seq; Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON SEQUENCE public.clients_id_seq TO postgres;
GRANT ALL ON SEQUENCE public.clients_id_seq TO anon;
GRANT ALL ON SEQUENCE public.clients_id_seq TO authenticated;
GRANT ALL ON SEQUENCE public.clients_id_seq TO service_role;


--
-- TOC entry 2909 (class 0 OID 0)
-- Dependencies: 252
-- Name: TABLE launches; Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON TABLE public.launches TO postgres;
GRANT ALL ON TABLE public.launches TO anon;
GRANT ALL ON TABLE public.launches TO authenticated;
GRANT ALL ON TABLE public.launches TO service_role;


--
-- TOC entry 2910 (class 0 OID 0)
-- Dependencies: 258
-- Name: SEQUENCE launches_id_seq; Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON SEQUENCE public.launches_id_seq TO postgres;
GRANT ALL ON SEQUENCE public.launches_id_seq TO anon;
GRANT ALL ON SEQUENCE public.launches_id_seq TO authenticated;
GRANT ALL ON SEQUENCE public.launches_id_seq TO service_role;


--
-- TOC entry 2911 (class 0 OID 0)
-- Dependencies: 255
-- Name: TABLE projects; Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON TABLE public.projects TO postgres;
GRANT ALL ON TABLE public.projects TO anon;
GRANT ALL ON TABLE public.projects TO authenticated;
GRANT ALL ON TABLE public.projects TO service_role;


--
-- TOC entry 2912 (class 0 OID 0)
-- Dependencies: 256
-- Name: SEQUENCE projects_id_seq; Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON SEQUENCE public.projects_id_seq TO postgres;
GRANT ALL ON SEQUENCE public.projects_id_seq TO anon;
GRANT ALL ON SEQUENCE public.projects_id_seq TO authenticated;
GRANT ALL ON SEQUENCE public.projects_id_seq TO service_role;


--
-- TOC entry 2913 (class 0 OID 0)
-- Dependencies: 257
-- Name: TABLE results; Type: ACL; Schema: public; Owner: supabase_admin
--

GRANT ALL ON TABLE public.results TO postgres;
GRANT ALL ON TABLE public.results TO anon;
GRANT ALL ON TABLE public.results TO authenticated;
GRANT ALL ON TABLE public.results TO service_role;


--
-- TOC entry 2331 (class 826 OID 16476)
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES  TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES  TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES  TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES  TO service_role;


--
-- TOC entry 2332 (class 826 OID 16477)
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: supabase_admin
--

ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON SEQUENCES  TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON SEQUENCES  TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON SEQUENCES  TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON SEQUENCES  TO service_role;


--
-- TOC entry 2330 (class 826 OID 16475)
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS  TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS  TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS  TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS  TO service_role;


--
-- TOC entry 2334 (class 826 OID 16479)
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: public; Owner: supabase_admin
--

ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON FUNCTIONS  TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON FUNCTIONS  TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON FUNCTIONS  TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON FUNCTIONS  TO service_role;


--
-- TOC entry 2329 (class 826 OID 16474)
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES  TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES  TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES  TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES  TO service_role;


--
-- TOC entry 2333 (class 826 OID 16478)
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: supabase_admin
--

ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON TABLES  TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON TABLES  TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON TABLES  TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON TABLES  TO service_role;
