export class User {
    public name: string;
}

export class RegistrationData {
    public name: string;
    public emailAddress: string;
    public phone: string;
    public role: string
    public riskPlanId: number;
    public userId: string;
    public password: string;
    public message: string;
    public tokenExpired: string;
    public confirmationToken: string;
    public isApiCall: boolean;
    public isEmailConfirmed: boolean;
}

export class RegisterResult {
    public isSuccessful: boolean;
    public tokenType: string;
    public accessToken: string;
    public expiresIn: string;
}