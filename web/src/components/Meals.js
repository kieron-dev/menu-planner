import React, { useState } from "react";
import { Form, Card } from "semantic-ui-react";
import Meal from "./Meal";
import { v4 } from "uuid";
import { useRecipes } from "./RecipeProvider";

export default function Meals() {
    const [mealRecipes, setMealRecipes] = useState([]);
    const [meals, setMeals] = useState([]);
    const { recipes, addRecipe } = useRecipes();

    const changeName = (_, { value }) => {
        setMealRecipes(value);
    };

    const handleAddition = (_, { value }) => {
        addRecipe(value).then((r) =>
            setMealRecipes(mealRecipes.concat(r.value))
        );
    };

    const addMeal = (e) => {
        e.preventDefault();
        const newMeals = meals.concat({ id: v4(), recipes: mealRecipes });
        setMeals(newMeals);
        setMealRecipes([]);
    };

    const removeMeal = (id) => {
        const newMeals = meals.filter((m) => m.id !== id);
        setMeals(newMeals);
    };

    return (
        <>
            <Form onSubmit={addMeal} style={{ marginBottom: "5ex" }}>
                <Form.Group>
                    <Form.Dropdown
                        placeholder="Choose meal recipes..."
                        value={mealRecipes}
                        options={[...recipes.values()]}
                        onChange={changeName}
                        onAddItem={handleAddition}
                        width="14"
                        search
                        selection
                        multiple
                        allowAdditions
                        additionPosition="bottom"
                    />
                    <Form.Button type="submit">Add Meal</Form.Button>
                </Form.Group>
            </Form>

            <Card.Group>
                {meals.map((m) => (
                    <Meal
                        id={m.id}
                        recipeIDs={m.recipes}
                        removeMeal={removeMeal}
                    />
                ))}
            </Card.Group>
        </>
    );
}
