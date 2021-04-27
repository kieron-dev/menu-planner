import React from "react";
import LoginWrapper from "./components/LoginWrapper";
import AuthProvider from "./components/AuthProvider";

import "semantic-ui-css/semantic.min.css";

function App() {
    return (
        <div className="App">
            <AuthProvider>
                <LoginWrapper />
            </AuthProvider>
        </div>
    );
}

export default App;
