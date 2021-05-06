import React, { createContext, useContext, useState, useEffect } from "react";
import { v4 } from "uuid";

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
        const newId = v4();
        const newRecipe = {
            key: name,
            text: name,
            value: newId,
        };
        const newRecipes = new Map(recipes).set(newId, newRecipe);
        setRecipes(newRecipes);

        return newRecipe;
    };

    return (
        <RecipeContext.Provider value={{ recipes, addRecipe }}>
            {children}
        </RecipeContext.Provider>
    );
}
