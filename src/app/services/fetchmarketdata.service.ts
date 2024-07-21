import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../environments/environments';
import { map, Observable } from 'rxjs';
import { IItemListings } from '../interfaces/itemlistings.interface';

@Injectable({
  providedIn: 'root',
})
export class FetchmarketdataService {
  constructor(private _http: HttpClient) {}

  getMarketData(): Observable<IItemListings[]> {
    return this._http
      .get<{ items: IItemListings[] }>(`${environment.apiUrl}/api/data`)
      .pipe(map((data) => data.items));
  }
}
