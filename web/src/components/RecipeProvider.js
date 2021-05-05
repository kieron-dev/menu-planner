import React, { createContext, useContext, useState } from "react";
import { v4 } from "uuid";

const RecipeContext = createContext();
export const useRecipes = () => useContext(RecipeContext);

export default function RecipeProvider({ children }) {
    const recipeMap = new Map();
    ["Spaghetti Bolognese", "Toad in the Hole", "Spicy Chicken Tagine"].forEach(
        (r) => {
            const id = v4();
            recipeMap.set(id, {
                key: r,
                text: r,
                value: id,
            });
        }
    );
    const [recipes, setRecipes] = useState(recipeMap);

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
