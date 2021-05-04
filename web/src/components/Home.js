import React from "react";
import TopMenu from "./TopMenu";
import { Container } from "semantic-ui-react";
import Meals from "./Meals";

export default function Home() {
    return (
        <>
            <TopMenu />
            <Container>
                <Meals />
            </Container>
        </>
    );
}
