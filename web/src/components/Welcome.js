import React from 'react';
import {
  Segment,
  Container,
  Header
} from 'semantic-ui-react';
import GoogleLogin from 'react-google-login';
import {useAuth} from './AuthProvider';

const Welcome = () => {

  const {authGoogle} = useAuth();

  const successGoogle = (resp) => {
    authGoogle(resp.tokenId);
  };

  const failureGoogle = (resp) => {
    console.log(resp);
  };

  return (
    <>
      <Segment
        inverted
        textAlign='center'
        vertical
        style={{minHeight: 700}}
      >
        <Container text>
          <Header
            inverted
            as="h1"
            style={{fontSize: "4em", fontWeight: "normal", marginTop: "3em", marginBottom: "0"}}
            >Menu-Planner</Header>
          <Header
            inverted
            as="h2"
            style={{fontSize: "1.7em", fontWeight: "normal", marginTop: "1.5em", marginBottom: "1.5em"}}
          >Sign in to experience the delights!</Header>
          <GoogleLogin
            clientId="176462381984-bfq3v9mc00v0ipvpebiaiide4l22dmoh.apps.googleusercontent.com"
            onSuccess={successGoogle}
            onFailure={failureGoogle}
          />

        </Container>
      </Segment>
    </>
  );
};

export default Welcome;
