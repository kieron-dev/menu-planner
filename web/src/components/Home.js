import React from "react";
import { useAuth } from "./AuthProvider";

const Home = () => {
    const { auth, logout } = useAuth();

    return (
        <>
            <h1>Hello {auth.name}!</h1>
            <button onClick={logout}>Logout</button>
        </>
    );
};

export default Home;
