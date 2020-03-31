import { Context } from 'koa';

export async function getTransaction (ctx : Context){
  
    if((ctx.params.playerAccountID.trim()).length > 0 && (ctx.params.transactionID.trim()).length > 0){
        var data : any = {
            "TransactionType": 0,
            "AccountNumber": 0,
            "AssetNumber": "abc",
            "VoucherNumber": 0,
            "Created": new Date(),
            "Updated": new Date(),
            "InstanceID": 0,
            "Status": 0,
            "State": 0,
			"Amount":"10",
        }  
        ctx.body = data
   }else{
       
   }
    
}

export async function transferFromWallet​ (ctx : Context){
    let validate = await paramValidation(ctx)
    if(validate === true){
        ctx.body = {}
    }else{
        ctx.body = "Invalid request"
    }
}

export async function transferFromETG ( ctx : Context){

    let validate = await paramValidation(ctx)
    if(validate === true){
        ctx.body = {"Amount":"10"}
    }else{
        ctx.body = "Invalid request"
    }
}

function paramValidation(ctx : Context){
   
    if((ctx.params.assetNumber.trim()).length > 0 && (ctx.params.playerAccountID.trim()).length > 0){
        return true
    }else{
        return false
    }
}