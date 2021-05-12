import React, { createContext, useContext, useState, useEffect } from "react";

const RecipeContext = createContext();
export const useRecipes = () => useContext(RecipeContext);

export default function RecipeProvider({ children }) {
    const [recipes, setRecipes] = useState([]);

    useEffect(() => {
        fetch(process.env.REACT_APP_API_URI + "/recipes", {
            credentials: "include",
            method: "GET",
        })
            .then((resp) => {
                if (!resp.ok) throw new Error(resp.statusText);
                return resp;
            })
            .then((r) => r.json())
            .then((json) => {
                const recipeMap = new Map();
                json.forEach((r) => {
                    recipeMap.set(r.id, {
                        key: r.name,
                        text: r.name,
                        value: r.id,
                    });
                });
                setRecipes(recipeMap);
            })
            .catch(console.error);
    }, []);

    const addRecipe = (name) => {
        return new Promise((resolve, reject) => {
            fetch(process.env.REACT_APP_API_URI + "/recipes", {
                credentials: "include",
                method: "POST",
                body: JSON.stringify({ name }),
            })
                .then((resp) => {
                    if (!resp.ok) {
                        reject(resp.statusText);
                        throw new Error(resp.statusText);
                    }

                    return resp;
                })
                .then((r) => r.json())
                .then((r) => {
                    const newRecipes = new Map(recipes);
                    const newRecipe = {
                        key: r.name,
                        text: r.name,
                        value: r.id,
                    };
                    newRecipes.set(newRecipe.value, newRecipe);
                    setRecipes(newRecipes);

                    resolve(newRecipe);
                });
        });
    };

    return (
        <RecipeContext.Provider value={{ recipes, addRecipe }}>
            {children}
        </RecipeContext.Provider>
    );
}
