import { TestBed } from '@angular/core/testing';

import { FetchmarketdataService } from './fetchmarketdata.service';

describe('FetchmarketdataService', () => {
  let service: FetchmarketdataService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FetchmarketdataService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
