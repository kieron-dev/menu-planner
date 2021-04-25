import React, { createContext, useState, useContext, useCallback } from "react";

const initAuth = {
    isAuthed: false,
    name: "Guest",
};

const AuthContext = createContext();
export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    const [auth, setAuth] = useState(initAuth);

    const setAuthenticated = useCallback(
        (name) => {
            const newAuth = { isAuthed: true, name: name };
            setAuth(newAuth);
        }, []);

    const logout = useCallback(
        () => {
            fetch(process.env.REACT_APP_API_URI + "/logout", {
                credentials: "include",
                method: 'POST',
            })
                .then(resp => {
                    if (!resp.ok) throw new Error(resp.statusText);
                    return resp;
                })
                .then(() => setAuth({ isAuthed: false }))
                .catch(console.error);
        }, []);

    const authGoogle = (token) => {
        fetch(process.env.REACT_APP_API_URI + "/authGoogle", {
            credentials: "include",
            method: "POST",
            body: JSON.stringify({ idToken: token }),
            headers: {
                "Content-Type": "application/json",
            },
        })
            .then((resp) => {
                if (!resp.ok) throw new Error(resp.statusText);
                return resp;
            })
            .then((data) => data.json())
            .then((data) => {
                setAuthenticated(data.name);
            })
            .catch((err) => console.error(err));
    };

    return (
        <AuthContext.Provider value={{ auth, authGoogle, setAuthenticated, logout }}>
            {children}
        </AuthContext.Provider>
    );
};
