import React from "react";
import { rest } from "msw";
import { setupServer } from "msw/node";
import { render, waitFor, screen } from "@testing-library/react";
import LoginWrapper from "./LoginWrapper";
import AuthProvider from "./AuthProvider";

import Home from "./Home";
import Welcome from "./Welcome";
jest.mock("./Home", () => () => <div data-testid="home">Home</div>);
jest.mock("./Welcome", () => () => <div data-testid="welcome">Welcome</div>);

const server = setupServer(
    rest.get(process.env.REACT_APP_API_URI + "/whoami", (_, res, ctx) => {
        return res(ctx.json({ name: "bob" }));
    })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("when authenticated", () => {
    it("shows the home page", async () => {
        render(
            <AuthProvider>
                <LoginWrapper />
            </AuthProvider>
        );

        await waitFor(() => screen.getByTestId("home"));
    });
});

describe("when not authenticated", () => {
    beforeEach(() => {
        server.use(
            rest.get(
                process.env.REACT_APP_API_URI + "/whoami",
                (_, res, ctx) => {
                    return res(ctx.status(401));
                }
            )
        );
    });

    it("shows the login page", async () => {
        render(
            <AuthProvider>
                <LoginWrapper />
            </AuthProvider>
        );

        await waitFor(() => screen.getByTestId("welcome"));
    });
});
