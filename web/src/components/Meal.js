import React from "react";
import { Card, Button } from "semantic-ui-react";

export default function Meal({ id, header, description, removeMeal }) {
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
                <Card.Header>{header}</Card.Header>
                <Card.Description>{description}</Card.Description>
            </Card.Content>
        </Card>
    );
}
