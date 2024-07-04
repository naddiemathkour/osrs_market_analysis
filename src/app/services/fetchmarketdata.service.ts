import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../environments/environments';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class FetchmarketdataService {
  constructor(private _http: HttpClient) {}

  getMarketData(): Observable<any> {
    return this._http.get(`${environment.apiUrl}api/data`);
  }
}
