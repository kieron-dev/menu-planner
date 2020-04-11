import React, {createContext, useState, useContext} from "react";

const initAuth = {
  isAuthed: false,
  name: "Guest"
};

const AuthContext = createContext();
export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({children}) => {
  const [auth, setAuth] = useState(initAuth);

  const authGoogle = (token) => {
    fetch(process.env.REACT_APP_API_URI + "/authGoogle", {
      method: "POST",
      body: JSON.stringify({tokenID: token}),
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(
      (data) => data.json()
    ).then(
      (json) => {
        const newAuth = {...auth, isAuthed: true, name: 'real name'};
        setAuth(newAuth);
      }
    ).catch(
      (err) => console.log(err)
    );
  };

  return (
    <AuthContext.Provider value={{auth, authGoogle}}>
      {children}
    </AuthContext.Provider>
  );
};
