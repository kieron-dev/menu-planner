import React from "react";
import TopMenu from "./TopMenu";
import { Container } from "semantic-ui-react";
import Meals from "./Meals";
import RecipeProvider from "./RecipeProvider";

export default function Home() {
    return (
        <>
            <TopMenu />
            <Container>
                <RecipeProvider>
                    <Meals />
                </RecipeProvider>
            </Container>
        </>
    );
}
