import React, { useState, useEffect } from "react";
import { useAuth } from "./AuthProvider";
import Home from "./Home";
import Welcome from "./Welcome";

const LoginWrapper = () => {
    const { auth, setAuthenticated, setUnauthenticated } = useAuth();
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetch(process.env.REACT_APP_API_URI + "/whoami", {
            credentials: "include",
            method: "GET",
        })
            .then((resp) => {
                if (!resp.ok) {
                    setLoading(false);
                    setUnauthenticated();
                    throw new Error("not-authed");
                }
                return resp;
            })
            .then((data) => data.json())
            .then((data) => setAuthenticated(data.name))
            .then(() => setLoading(false))
            .catch(() => {});
    }, [setUnauthenticated, setAuthenticated]);

    if (loading) return <p>Loading...</p>;

    return auth.isAuthed ? <Home /> : <Welcome />;
};

export default LoginWrapper;
