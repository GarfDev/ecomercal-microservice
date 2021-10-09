import { Injectable, Request } from '@nestjs/common';

@Injectable()
export class AppService {
  getData(req: Request): { message: string } {
    console.log(req?.headers)
    return { message: 'Welcome to user-service!' };
  }
}
