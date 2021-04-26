import React from 'react';
import { rest } from 'msw';
import { setupServer } from 'msw/node';
import { render, waitFor, screen } from '@testing-library/react';
import App from '../App';

const server = setupServer(
    rest.get(process.env.REACT_APP_API_URI + '/whoami', (_, res, ctx) => {
        return res(ctx.json({ name: 'bob' }));
    })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe('when authenticated', () => {
    it('shows the home page', async () => {
        render(<App />);

        await waitFor(() => screen.getByRole('button'));
        expect(screen.getByRole('button')).toHaveTextContent(/logout/i);
    });
});

describe('when not authenticated', () => {
    beforeEach(() => {
        server.use(
            rest.get(process.env.REACT_APP_API_URI + '/whoami', (_, res, ctx) => {
                return res(ctx.status(401));
            })
        );
    });

    it('shows the login page', async () => {
        render(<App />);

        await waitFor(() => screen.getByRole('button'));
        expect(screen.getByRole('button')).toHaveTextContent(/google/i);
    });
});
