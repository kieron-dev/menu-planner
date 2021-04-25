import React from 'react';
import { render, cleanup, waitForElement } from '@testing-library/react';
import App from '../App';

describe('login wrapper', () => {
    var fakeFetch;

    beforeEach(() => {
        fakeFetch = jest.spyOn(window, 'fetch');
    });

    afterEach(() => {
        cleanup()
    });

    it('calls the /whoami endpoint', async () => {
        fakeFetch.mockReturnValue(new Promise(f => f));
        render(<App />);

        const prefix = process.env.REACT_APP_API_URI;
        expect(fakeFetch).toHaveBeenCalledWith(
            `${prefix}/whoami`,
            {
                credentials: 'include',
                method: 'GET',
            }
        );

    });

    describe('when logged in', () => {
        it('shows the home page', async () => {
            fakeFetch.mockResolvedValueOnce({
                ok: true,
                json: async () => ({
                    "name": "bob"
                })
            })
            const { getByText } = render(<App />);
            const logoutBtn = await waitForElement(() => getByText(/Logout/i));
            expect(logoutBtn).toBeInTheDocument();
        });
    });

    describe('when not logged in', () => {
        it('shows login page', async () => {
            fakeFetch.mockResolvedValue({ ok: false });
            const { getByText } = render(<App />);
            const googleLogin = await waitForElement(() => getByText(/Sign in with Google/i));
            expect(googleLogin).toBeInTheDocument();
        });
    });
});
