Read this article: https://miparnisariblog.wordpress.com/2024/07/01/solving-identity-management-in-modern-applications-book-notes/

Events in the life of a user:

    Provisioning
    Authorization (granting of privileges)
    Authentication
    Access policy enforcement
    Sessions
    Single sign-on (a user authenticates once and accesses multiple applications, without having to reauthenticate because they all trust the same identity provider)
    Stronger authentication (supported by both OIDC and SAML 2)
    Logout (may involve terminating more than one session)
    Account management
    Deprovisioning

Protocols

    SAML 2: solves the problem of single sign-on across domains. It’s complex to implement but it’s widely used.
    Ws-Fed: also solves single sign-on.
    OAuth 2: allows a user to authorize one application (the “client”) to send a request to an API (the “resource server”) on the user’s behalf, without needing their credentials. It solves authorization. Current version is OAuth 2.0 but 2.1 is on the works.
    OIDC: layer on top of OAuth 2 that provides information to applications about the identity of authenticated users. Solves authentication.
    SCIM (System for Cross-Domain Identity Management): solves the problem of sending and updating identity information from one domain to another. It provides a standard REST API for one system to send requests to another for adding or modifying user and group records.

OAuth 2

    Example: user is writing documents in WriteAPaper.com using his own documents stored in documents.com. So WriteAPaper.com requests access to those documents with the user’s consent.
    An authorization server must implement the /authorize endpoint and the /oauth/token endpoint, which accepts a code.
    Three grant types / flows:
        Authorization code with PKCE (Proof Key for Code Exchange): the code ensures that the application that requested the authorization code is the same application that uses said code to get an access token. A random string called code verifier is used to create a code challenge. The authorization server must remember the code challenge. Then, when it receives a call to /oauth/token with the code verifier, it uses the same hash function to confirm that the code challenge is the same.
        Client credentials: requires no end-user interaction because the user doesn’t own the resource.
        Refresh token: makes sense only for the authorization code flow, to avoid having to request user consent every time.
        Client device (not part of the spec): used for IoT, for example, a digital picture frame that requests permission to show pictures from another site. The user grants permission from a secondary device, like a phone.
        Implicit flow: not recommended because it exposes access tokens in URL fragments which could be stored in the browser’s history.

Authorization Code with PKCE
OIDC

    An authorization server can implement a /userinfo endpoint which returns claims with a JSON object. This is useful if the claims are too large to fit in an ID token, which is encoded as JWT. The /userinfo endpoint can only be accessed with a valid access token.
    Three grant types / flows:
        Authorization code flow: similar to OAuth 2’s flow.
        Implicit flow: similar to OAuth2’s flow except that there is no risk of exposing an access token because the authorization server only returns an ID token in the form of a JWT.
        Hybrid flow: mixture of the previous two. It is designed for applications with both a secure back end and a front-end.

SAML 2

    Provides cross-domain single sign-on and identity federation (i.e. a way for an application and an identity provider to use a common shared identifier for a user).
    The identity provider must implement a /sso endpoint where SAML requests can be received. The application must implement an /acs endpoint where the SAML response can be received.

Authorization

Authorization is the granting of privileges to access a resource. Access policy enforcement is done when a user requests a resource and a check is made to see if they can have access.
Advertisement
Privacy Settings

Three levels at which authorization may be specified and applied:

    At the application or API
    What functions can the user call in the application or API
    What data the user can access or operate on

Models for specifying authorization:

    Access control lists
    RBAC
    ABAC
    ReBAC
    JWT-secured Authorization Requests (JARs), Rich Authorization Requests (RARs), or Pushed Authorization Requests (PARs).

Compliance

    GDPR describes a legal basis for processing personal data. It applies to any product that processes personal information of EU residents, regardless of where the information is held.
    Payment card industry must comply with PCI DSS.
    HIPAA and HITECH is required for the healthcare industry in the United States.
    FedRAMP applies to companies providing services to US government agencies.
    SOC2 (Service Organization Control) is a set of controls against which a company is audited related to security, privacy, confidentiality, integrity and availability.
