import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../environments/environments';
import { catchError, map, Observable, of, switchMap, timer } from 'rxjs';
import { IItemListings } from '../interfaces/itemlistings.interface';

@Injectable({
  providedIn: 'root',
})
export class FetchmarketdataService {
  constructor(private _http: HttpClient) {}

  getMarketData(): Observable<IItemListings[]> {
    return timer(0, 0.5 * 60 * 1000).pipe(
      switchMap(() =>
        this._http
          .get<{ items: IItemListings[] }>(`${environment.apiUrl}api/data`)
          .pipe(
            map((data) => data.items),
            catchError((error) => {
              console.error('Error fetching market data:', error);
              return of([] as IItemListings[]);
            })
          )
      )
    );
  }
}
