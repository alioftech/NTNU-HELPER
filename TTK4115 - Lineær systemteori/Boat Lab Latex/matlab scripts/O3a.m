%%From earlier excercies 
K=0.1734;
T=84.3920;
%% 3.1a
omega_c=0.1; %cutoff frequency [rad/s]
PM=50/180*pi;%Phase margin [rad]
T_d=T;        %chosen such that it cancels the TF time constant
%Make transfer function for controller
T_f=1/(tan(PM)*omega_c);    
K_pd=sqrt((T_f^2*omega_c^4+omega_c^2)/K^2);
num_controller=[K_pd*T_d,K_pd];
den_controller=[T_f,1];

H_pd=tf(num_controller,den_controller); %make transfer function for controller

%%Make transfer function for plant
H_ship=tf([K],[T 1 0]);                 %transfer function for plant

%open-loop system
H_ol=H_pd*H_ship;                       %Open loop transfer function

%draw a bode idagram
figure
bode(H_ol);
grid ;
title('Bode plot for the open-loop system H_ol(s)');



