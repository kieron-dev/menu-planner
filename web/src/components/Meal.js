import React from "react";
import { Card, Button } from "semantic-ui-react";
import { useRecipes } from "./RecipeProvider";

export default function Meal({ id, recipeIDs, removeMeal }) {
    const { recipes } = useRecipes();

    const mainRecipe = recipes.get(recipeIDs[0]);

    return (
        <Card id={id}>
            <Card.Content>
                <Button
                    basic
                    floated="right"
                    size="mini"
                    icon="close"
                    onClick={() => removeMeal(id)}
                />
                <Card.Header>{mainRecipe.text}</Card.Header>
                {recipeIDs.map((id) => (
                    <Card.Description>{recipes.get(id).text}</Card.Description>
                ))}
            </Card.Content>
        </Card>
    );
}
