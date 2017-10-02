import { Injectable } from '@angular/core';
import { Device, Capability } from '../models/device-models';
import { Headers, Http } from '@angular/http';
import { AuthService } from './auth.service';
import { Router } from '@angular/router';
import { environment } from '../../environments/environment';
import  * as HttpStatus from 'http-status-codes'

@Injectable()
export class DeviceService {

    constructor(private http: Http, private authService: AuthService,
        private router: Router) {
    }

    public getDevices(userId: number): Promise<Array<Device>> {
        let headers = this.authService.initAuthHeaders();
        let url = `/api/${userId}/devices`;        

        return this.http.get(url, { headers: headers })
            .toPromise()
            .then(response => {
                let res = response.json() as Array<Device>;
                // res.forEach(x => x.expidationDate !== undefined && x.expidationDate != null ? x.expidationDate = new Date(x.expidationDate as any) : null);
                // res = res.filter(x => x.status !== MeasureStatus.dropped);
                return res;
            })
            .catch(e => {
                if (e.status === HttpStatus.UNAUTHORIZED) {
                    this.authService.logout();
                    this.router.navigate([environment.urls.login]);
                } else {
                    throw e;
                }
            });
    }

    public getDevice(id: number): Promise<Device> {
        let headers = this.authService.initAuthHeaders();
        let url = `/api/measure/${id}`;

        return this.http.get(url, { headers: headers })
            .toPromise()
            .then(response => {
                let res = response.json() as Device;
                // res.expidationDate = res.expidationDate !== undefined && res.expidationDate != null ?  new Date(res.expidationDate) as any : null;
                return res;
            })
            .catch(e => {                
                if (e.status === 401) {
                    this.authService.logout();
                    this.router.navigate([environment.urls.login]);
                } else {
                    throw e;
                }
            });
    }
}