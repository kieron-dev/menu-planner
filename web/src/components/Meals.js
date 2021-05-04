import React, { useState } from "react";
import { Form, Card, Button, Dropdown } from "semantic-ui-react";
import Meal from "./Meal";
import { v4 } from "uuid";

export default function Meals() {
    const [mealRecipes, setMealRecipes] = useState([]);
    const [recipes, setRecipes] = useState(
        ["Spaghetti Bolognese", "Toad in the Hole", "Spicy Chicken Tagine"].map(
            (r) => ({
                key: r,
                text: r,
                value: v4(),
            })
        )
    );

    const [meals, setMeals] = useState([]);

    const changeName = (_, { value }) => {
        console.log("value", value);
        setMealRecipes(value);
    };

    const handleAddition = (_, { value }) => {
        const newId = v4();
        const newRecipes = recipes.concat({
            key: value,
            text: value,
            value: newId,
        });
        setRecipes(newRecipes);
        setMealRecipes(mealRecipes.concat(newId));
    };

    const addMeal = (e) => {
        e.preventDefault();
        const newMeals = meals.concat({ id: v4(), name: mealRecipes });
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
                <Form.Field>
                    <Dropdown
                        placeholder="Add recipes"
                        value={mealRecipes}
                        options={recipes}
                        onChange={changeName}
                        onAddItem={handleAddition}
                        fluid
                        search
                        selection
                        multiple
                        allowAdditions
                        additionPosition="bottom"
                    />
                </Form.Field>
                <Button type="submit">Add Meal</Button>
            </Form>

            <Card.Group>
                {meals.map((m) => (
                    <Meal
                        id={m.id}
                        header={m.name}
                        description={m.description}
                        removeMeal={removeMeal}
                    />
                ))}
            </Card.Group>
        </>
    );
}
