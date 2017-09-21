import { AuthRequestResult } from './../models/request-result';
import { RegistrationData, User } from './../models/auth-models';
import { PersistenceService } from './persistence.service'
import { ConfigService} from './config.service'
import { Injectable } from "@angular/core";
import { Headers, Http } from "@angular/http";
import "rxjs/add/operator/toPromise";
import { Router } from '@angular/router';

@Injectable()
export class AuthService {
    private tokeyKey = "token";
    private token: string;

    constructor(
        private http: Http,        
        private persistance: PersistenceService,
        private router: Router
    ) { }

    public login(userName: string, password: string): Promise<AuthRequestResult> {
        let body = { Username: userName, Password: password };

        return this.http.post('/api/TokenAuth', body).toPromise()
            .then(response => {
                let result = response.json() as AuthRequestResult;
                if (result.State == 1) {
                    let json = result.Data as any;
                    this.persistance.accessToken = json.accessToken;
                    this.persistance.currentUser = json.user;
                }
                return result;
            })
            .catch(this.handleError);
    }

    public logout(): void {
        this.token == null;
        this.persistance.accessToken = null;
        this.persistance.planId = null;
        this.persistance.currentUser = null;
        this.persistance.showCompltetedMeasures = true;
    }

    public checkLogin(): boolean {
        let token = this.persistance.accessToken;
        return token != null;
    }

    public getUserInfo(): Promise<AuthRequestResult> {
        return this.authGet('/api/TokenAuth');
    }

    public authPost(url: string, body: any): Promise<AuthRequestResult> {
        let headers = this.initAuthHeaders();
        return this.http.post(url, body, { headers: headers }).toPromise()
            .then(response => {
                let res = response.json() as AuthRequestResult
                return res;
            })
            .catch(this.handleError);
    }

    public authGet(url): Promise<AuthRequestResult> {
        let headers = this.initAuthHeaders();
        return this.http.get(url, { headers: headers }).toPromise()
            .then(response => {
                let res = response.json() as AuthRequestResult;
                return res;
            })
            .catch(this.handleError);
    }

    public createUser(userId: string, token: string): Promise<void> {
        let url = '/api/user';
        let body = {
            UserId: userId,
            Token: token
        };
        return this.http.post(url, body).toPromise()
            .then(response => { })
            .catch(this.handleError);
    }

    public getRegistrationData(userId: string, token: string, planId: number): Promise<RegistrationData> {
        let url = '/api/invitation';
        let body = {
            userId: userId,
            token: token,
            planId: planId
        };   

        return this.http.post(url, body).toPromise()
            .then(response => {
                let res = response.json() as RegistrationData;
                return res
            })
            .catch();
    }

    //return any for now, create a view model later
    public register(data: RegistrationData): Promise<any> {
        let url = '/api/invitation/register';
        data.isApiCall = true;

        return this.http.post(url, data).toPromise()
            .then(response => {
                let res = response.json() as any;
                this.persistance.accessToken = res.accessToken;
                this.persistance.currentUser = res.user;

                return res
            })
            .catch();
    }

    private getLocalToken(): string {
        this.token = this.persistance.accessToken;
        return this.token;
    }

    public initAuthHeaders(redirectIfNoToken: boolean = true): Headers {
        let token = this.getLocalToken();
        if (token == null) {
            if (redirectIfNoToken) {
                this.logout();
                this.router.navigate['/login']
            } else {
                throw 'No token';
            }
        }

        var headers = new Headers();
        headers.append('Authorization', 'Bearer ' + token);

        return headers;
    }

    private handleError(error: any): Promise<any> {
        console.error('An error occurred', error);
        return Promise.reject(error.message || error);
    }
}