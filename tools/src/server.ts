import koa, { Context } from 'koa';
import cors from '@koa/cors'
import bodyParse from 'koa-body';
import Router from 'koa-router';
import * as api from './app/smib/smib'

const app = new koa();
const router = new Router();
var port : number = 8011;
console.log(__dirname);

app.use(cors());
app.use(bodyParse());

router.get('/api/SMIB/Transaction/:playerAccountID/:transactionID', api.getTransaction)

router.post('/api/SMIB/TransferFromWallet/:assetNumber/:playerAccountID', validate, api.transferFromWallet​);

router.post('/api/SMIB/TransferFromETG/:assetNumber/:playerAccountID', validate, api.transferFromETG)

function validate (ctx : Context, next : any) {
    if(ctx.request.body){
        next()
    }else{
        ctx.status = 400
    }
}

app.use(router.routes());

app.listen(port, function(){
    console.log(`server runs on ${port}`);
})