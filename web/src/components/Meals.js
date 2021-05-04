import React, { useState } from "react";
import { Form, Card, Button } from "semantic-ui-react";
import Meal from "./Meal";
import v4 from "uuid";

export default function Meals() {
    const [name, setName] = useState("");

    const [meals, setMeals] = useState(
        ["spag bol", "lasagne"].map((n) => ({
            id: v4(),
            name: n,
            description: "blah",
        }))
    );

    const changeName = (e) => {
        setName(e.target.value);
    };

    const addMeal = (e) => {
        e.preventDefault();
        const newMeals = meals.concat({ id: v4(), name: name });
        setMeals(newMeals);
        setName("");
    };

    const removeMeal = (id) => {
        const newMeals = meals.filter((m) => m.id !== id);
        setMeals(newMeals);
    };

    return (
        <>
            <Form onSubmit={addMeal} style={{ marginBottom: "5ex" }}>
                <Form.Input
                    label="Recipe"
                    placeholder="Recipe"
                    value={name}
                    onChange={changeName}
                />
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
